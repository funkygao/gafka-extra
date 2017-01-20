package main

import (
	"fmt"
	"sort"
	"time"

	"github.com/funkygao/golib/bjtime"
	"github.com/funkygao/termui"
	"github.com/mattn/go-runewidth"
	"github.com/nsf/termbox-go"
	"github.com/pmylund/sortutil"
)

var (
	w, h               int
	page               int
	pageSize           int
	selectedRow        = 0
	selectedCommit     = 0
	currentWebHook     *Webhook
	mainView           bool
	detailView         = false
	dashboardView      = false
	userSummaryView    = false
	projectSummaryView = false
)

const (
	coldef = termbox.ColorDefault
)

func refreshSize() {
	w, h = termbox.Size()
	pageSize = h - 1 // footer=1
}

func redrawAll() {
	mainView = true

	termbox.Clear(coldef, coldef)
	refreshSize()

	x, y := 0, 0
	lock.Lock()
	evts := reorderEvents()
	for i := page * pageSize; i < len(evts) && i < (page+1)*pageSize; i++ {
		evt := evts[i]
		drawEvent(x, y, evt)
		y++
	}

	drawFooter()
	lock.Unlock()

	termbox.Flush()
}

func drawDashboardByUser() {
	lock.Lock()
	commitByUsers := make(map[string]int)
	for _, evt := range events {
		if hook, ok := evt.(*Webhook); ok {
			for _, c := range hook.Commits {
				key := fmt.Sprintf("%s %s", c.Author.Name, c.Author.Email)
				if _, present := commitByUsers[key]; !present {
					commitByUsers[key] = 1
				} else {
					commitByUsers[key] += 1
				}
			}
		}
	}
	lock.Unlock()

	type user struct {
		name    string
		commits int
	}
	users := make([]user, 0, len(commitByUsers))
	for name, n := range commitByUsers {
		users = append(users, user{name: name, commits: n})
	}
	sortutil.DescByField(users, "commits")

	termbox.Clear(coldef, coldef)
	y := 0
	for i, u := range users {
		row := fmt.Sprintf("%s: %5d", wideStr(u.name, 60), u.commits)
		drawRow(row, y, coldef, coldef)
		y++

		if i >= h {
			break
		}
	}

	termbox.Flush()
}

func drawDashboardByProject() {
	lock.Lock()
	commitByProjects := make(map[string]int)
	for _, evt := range events {
		if hook, ok := evt.(*Webhook); ok {
			if _, present := commitByProjects[hook.Repository.Homepage]; !present {
				commitByProjects[hook.Repository.Homepage] = hook.Total_commits_count
			} else {
				commitByProjects[hook.Repository.Homepage] += hook.Total_commits_count
			}
		}
	}
	lock.Unlock()

	type project struct {
		name    string
		commits int
	}
	projects := make([]project, 0, len(commitByProjects))
	for name, n := range commitByProjects {
		projects = append(projects, project{name: name, commits: n})
	}
	sortutil.DescByField(projects, "commits")

	termbox.Clear(coldef, coldef)
	y := 0
	for i, p := range projects {
		row := fmt.Sprintf("%80s: %5d", p.name, p.commits)
		drawRow(row, y, coldef, coldef)
		y++

		if i >= h {
			break
		}
	}

	termbox.Flush()
}

func drawDashboardByHour() {
	lock.Lock()
	commitByHours := make(map[int]int)
	for _, evt := range events {
		if hook, ok := evt.(*Webhook); ok {
			for _, cmt := range hook.Commits {
				t, _ := time.Parse(time.RFC3339, cmt.Timestamp)
				if _, present := commitByHours[t.Hour()]; present {
					commitByHours[t.Hour()]++
				} else {
					commitByHours[t.Hour()] = 1
				}
			}
		}
	}
	lock.Unlock()

	sortedHours := make([]int, 0, len(commitByHours))
	for hour := range commitByHours {
		sortedHours = append(sortedHours, hour)
	}
	sort.Ints(sortedHours)

	termbox.Clear(coldef, coldef)

	termui.UseTheme("helloworld")
	bc := termui.NewBarChart()
	data := []int{}
	bclabels := []string{}
	for _, hour := range sortedHours {
		data = append(data, commitByHours[hour])
		bclabels = append(bclabels, fmt.Sprintf("%02d", hour))
	}
	bc.Border.Label = "Commit by hour"
	bc.Data = data
	bc.Width = w
	bc.Height = h
	bc.DataLabels = bclabels
	bc.TextColor = termui.ColorGreen
	bc.BarColor = termui.ColorRed
	bc.NumColor = termui.ColorYellow

	termui.Render(bc)
}

func drawDetail() {
	evts := reorderEvents()
	evt := evts[selectedRow]
	if webhook, ok := evt.(*Webhook); !ok {
		return
	} else {
		currentWebHook = webhook
	}

	termbox.Clear(coldef, coldef)

	hook := evt.(*Webhook)

	y := 1
	fg, bg := coldef, coldef
	row := fmt.Sprintf("%7s: %s (%s)",
		"Repo",
		hook.Repository.Name, hook.Repository.Description)
	drawRow(row, y, fg, bg)
	y++

	row = fmt.Sprintf("%7s: %s", "Home", hook.Repository.Homepage)
	drawRow(row, y, fg, bg)
	y++

	row = fmt.Sprintf("%7s: %s    %7s: %s    %7s: %d    %7s: %d",
		"Ref", hook.Ref,
		"User", hook.User_name,
		"Commits", hook.Total_commits_count, "Project", hook.Project_id)
	drawRow(row, y, fg, bg)
	y += 2

	for i := len(hook.Commits) - 1; i >= 0; i-- {
		c := hook.Commits[i]
		row = fmt.Sprintf("%14s %20s %s", since(c.Timestamp), c.Author.Email, c.Message)
		drawRow(row, y, termbox.ColorGreen, bg)
		y++

		row = fmt.Sprintf("%s", c.Url)
		if selectedCommit == len(hook.Commits)-1-i {
			// selected commit row
			drawRow(row, y, termbox.ColorBlack, termbox.ColorYellow)
		} else {
			drawRow(row, y, fg, bg)
		}

		y++
	}

	termbox.Flush()
}

func drawRow(row string, y int, fg, bg termbox.Attribute) {
	x := 0
	for _, r := range row {
		termbox.SetCell(x, y, r, fg, bg)
		// wide string must be considered
		w := runewidth.RuneWidth(r)
		if w == 0 || (w == 2 && runewidth.IsAmbiguousWidth(r)) {
			w = 1
		}
		x += w
	}

	if isSelectedRow(y) {
		// highlight the whole line with spaces
		for i := x; i < w; i++ {
			termbox.SetCell(i, y, ' ', fg, bg)
		}
	}
}

func isSelectedRow(y int) bool {
	return y == selectedRow-(pageSize*page)
}

func drawEvent(x, y int, evt interface{}) {
	fg_col, bg_col := coldef, coldef
	if isSelectedRow(y) {
		fg_col = termbox.ColorBlack
		bg_col = termbox.ColorYellow
	}

	var row string
	switch hook := evt.(type) {
	case *Webhook:
		if len(hook.Commits) == 0 {
			row = fmt.Sprintf("%14s %s %-25s",
				bjtime.TimeToString(hook.ctime),
				wideStr(hook.User_name, 20),
				hook.Repository.Name)
		} else {
			commit := hook.Commits[len(hook.Commits)-1] // the most recent commit
			row = fmt.Sprintf("%14s %s %-25s %s",
				since(commit.Timestamp),
				wideStr(hook.User_name, 20),
				hook.Repository.Name,
				commit.Message)
		}

	case *SystemHookProjectCreate:
		fg_col = termbox.ColorRed
		row = fmt.Sprintf("%14s %20s created project(%s)",
			since(hook.Created_at),
			hook.Owner_name,
			hook.Name)

	case *SystemHookProjectDestroy:
		fg_col = termbox.ColorRed
		row = fmt.Sprintf("%14s %20s destroy project(%s)",
			since(hook.Created_at),
			hook.Owner_name,
			hook.Path_with_namespace)

	case *SystemHookGroupCreate:
		fg_col = termbox.ColorRed
		row = fmt.Sprintf("%14s %20s created group(%s)",
			since(hook.Created_at),
			hook.Owner_name,
			hook.Name)

	case *SystemHookGroupDestroy:
		fg_col = termbox.ColorRed
		row = fmt.Sprintf("%14s %20s destroy group(%s)",
			since(hook.Created_at),
			hook.Owner_name,
			hook.Name)

	case *SystemHookUserCreate:
		fg_col = termbox.ColorRed
		row = fmt.Sprintf("%14s %20s %s signup",
			since(hook.Created_at),
			hook.Name,
			hook.Email)

	case *SystemHookUserAddToGroup:
		fg_col = termbox.ColorRed
		row = fmt.Sprintf("%14s %20s join group(%s)",
			since(hook.Created_at),
			hook.User_name,
			hook.Group_name)

	case *SystemHookUserRemovedFromGroup:
		fg_col = termbox.ColorRed
		row = fmt.Sprintf("%14s %20s kicked from group(%s)",
			since(hook.Created_at),
			hook.User_name,
			hook.Group_name)

	case *SystemHookUserAddToTeam:
		fg_col = termbox.ColorRed
		row = fmt.Sprintf("%14s %20s join project(%s)",
			since(hook.Created_at),
			hook.User_name,
			hook.Project_name)

	case *SystemHookUserRemovedFromTeam:
		fg_col = termbox.ColorRed
		row = fmt.Sprintf("%14s %20s kicked from project(%s)",
			since(hook.Created_at),
			hook.User_name,
			hook.Project_name)

	case *SystemHookKeyCreate:
		fg_col = termbox.ColorRed
		row = fmt.Sprintf("%14s %20s create ssh key",
			since(hook.Created_at),
			hook.Username)

	case *SystemHookKeyDesctroy:
		fg_col = termbox.ColorRed
		row = fmt.Sprintf("%14s %20s destroy ssh key",
			since(hook.Created_at),
			hook.Username)

	case *SystemHookUnknown:
		fg_col = termbox.ColorMagenta
		row = fmt.Sprintf("%s", hook.Evt)
	}

	drawRow(row, y, fg_col, bg_col)

}

func drawSplash() {
	refreshSize()
	row := "loading gitlab events..."
	x, y := w/2-len(row)/2, h/2+1
	for i, c := range row {
		termbox.SetCell(x+i, y, c, termbox.ColorGreen, coldef)
	}
	termbox.Flush()
}

func drawNotify(row string) {
	refreshSize()
	termbox.Clear(coldef, coldef)
	x, y := w/2-len(row)/2, h/2+1
	for i, c := range row {
		termbox.SetCell(x+i, y, c, termbox.ColorGreen, coldef)
	}
	termbox.Flush()

	displayNotify(row, "Glass")
}

func drawFooter() {
	s := calculateStats()
	help := "q:Close d:Detail G:End K:Top Space:PageDown b:PageUp"
	stats := fmt.Sprintf("[events:%d/%d-%d page:%d/%d repo:%d staff:%d commit:%d]",
		selectedRow,
		loadedN,
		s.eventN,
		page,
		s.eventN/pageSize,
		s.repoN,
		s.staffN,
		s.commitN)
	footerText := help
	for i := 0; i < w-len(help)-len(stats); i++ {
		footerText += " "
	}
	footerText += stats
	for i := 0; i < w; i++ {
		termbox.SetCell(i, h-1, ' ', coldef, termbox.ColorBlue)
	}
	for i, c := range footerText {
		termbox.SetCell(i, h-1, c, coldef, termbox.ColorBlue)
	}
}
