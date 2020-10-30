package georaft



type Log struct {
    Term    int32         "term when entry was received by leader"
    Command interface{} "command for state machine,"
} 

type ApplyMsg struct {
    CommandValid bool
    Command      interface{}
    CommandIndex int
}