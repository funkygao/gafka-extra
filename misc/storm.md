No intermediate message brokers!
Higher level abstraction than message passing
“Just works”

Spouts can either be reliable or unreliable
Every topology has a "message timeout" associated with it
The actual topology has implicit streams and an implicit "acker" bolt added to manage the acking framework 

public interface IComponent extends Serializable {
    void declareOutputFields(OutputFieldsDeclarer outputFieldsDeclarer);

    Map<String, Object> getComponentConfiguration();
}

public interface ISpout extends Serializable {
    void open(Map conf, TopologyContext ctx, SpoutOutputCollector collector);
    void close();

    void activate();
    void deactivate();

    void nextTuple();

    // Acknowledges that a specific tuple is processed
    void ack(Object msgId);

    // Specifies that a specific tuple is not fully processed and to be reprocessed
    void fail(Object msgId);
}

public interface IRichSpout extends ISpout, IComponent {
}

public interface IBolt extends Serializable {
    void prepare(Map conf, TopologyContext ctx, OutputCollector collector);

    void execute(Tuple tup);

    void cleanup();
}

public interface IRichBolt extends IBolt, IComponent {
}
