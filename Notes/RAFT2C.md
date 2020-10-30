#### LAB 2C



> 1. Complete the functions `persist()` and `readPersist()` in `raft.go` by adding code to save and restore persistent state. You will need to encode (or "serialize") the state as an array of bytes in order to pass it to the `Persister`. Use the `labgob` encoder; see the comments in `persist()` and `readPersist()`. `labgob` is like Go's `gob` encoder but prints error messages if you try to encode structures with lower-case field names.
>
> 2. Insert calls to `persist()` at the points where your implementation changes persistent state. Once you've done this, you should pass the remaining tests.

由Task中的内容我们可以知道，Lab 2C的主要任务是实现persist()和readpersist()功能，并且将其放到合适的位置，完成raft的持久化和故障恢复操作操作。

#### 1. 实现persist()和readpersist()

>  save Raft's persistent state to stable storage, where it can later be retrieved after a crash and restart. see paper's Figure 2 for a description of what should be persistent.

在代码中有一段提示：应该根据Figure 2 中的描述那些raft的状态应该被persist。不过按照我的理解我只是对loglastAplied lastAplied commitIndex isVoted做了持久化。至于这样做是否存在有些状态多余或者缺少一些重要状态我没有深入思考。按照代码中原有的例子我们可以很快完成persist和readpersist两个函数。

~~~go
func (rf *Raft) persist() {
	 w := new(bytes.Buffer)
	 e := labgob.NewEncoder(w)
	 e.Encode(rf.lastAplied)
	 e.Encode(rf.log)
	 e.Encode(rf.lastAplied)
	 e.Encode(rf.commitIndex)
	 e.Encode(rf.isVoted)
	 data := w.Bytes()
	 rf.persister.SaveRaftState(data)
}
~~~

在完成了persist和readpersist函数以后，我们接下来要解决的是这两个函数应该在什么时候使用。readPersist在raft启动的时候使用就可以了，但是persist在什么时候应该调用呢？因为我对raft的理解感觉还不够深入，所以在persist()中包含的状态，改变以后我就调用一次persist()。这样做的好处是不会发生raft状态的丢失，但是因为做persist是需要io时间的，这使得写压力比较大的时候persist()的平凡调用会影响raft的性能。其实我觉得这个问题是挺值得思考的，不过目前还没有能力做这方面的优化，所以先在这里给自己挖个坑，以后来来填。

#### 1. 实现对确定Nextindex过程的优化

在Lab2B的实现中，我们find nextindex的方式非常的简单粗暴：从最新的log一个一个的试直到找到了正确的nextindex使得follower.log[nextindex-1].term == leader.log[index].term。但是这样做的话意味着每次通信都只能检测一位，当follower.log存在很长一段与leader冲突的term时，整个的效率就会变得非常低。

> For example, when rejecting an AppendEntries request, the follower can include the term of the conflicting entry and the first index it stores for that term. With this information, the leader can decrement nextIndex to bypass all of the conflicting entries in that term; one AppendEntries RPC will be required for each term with conflicting entries, rather than one RPC per entry. In practice, we doubt this optimization is necessary, since failures happen infrequently and it is unlikely that there will be many inconsistent entries.

按照实验中的提示我们很快就能找到论文中描述的一种优化的方法，虽然作怀疑在实际中这样做的价值，但是**在Lab2C中如果不先优化这一步提高find nextindex的效率，后面的测试节点是过不了的。**

优化后的代码：

~~~go
// 在AppendEntryReply中添加ConflictIndex和ConflictTerm
type AppendEntryReply struct{
	Term int
	Success bool 
	ConflictIndex int   
    ConflictTerm  int
}


func (rf *Raft) AppendEntries(args *AppendEntry, reply *AppendEntryReply) {
	......
    ......
    //// 在AppendEntries中添加确定ConflictIndex的操作
    reply.ConflictIndex = 1
	reply.ConflictTerm = -1
	prevLogTerm := -1
	logSize := len(rf.log)
	if args.PrevLogIndex >= 0 && args.PrevLogIndex < len(rf.log){
		prevLogTerm = rf.log[args.PrevLogIndex].Term
	}
	if prevLogTerm != args.PrevLogTerm {
		reply.ConflictIndex = logSize
		if prevLogTerm == -1 {
			
        } else { 
            reply.ConflictTerm = prevLogTerm 
			i := 0
            for ; i < logSize; i++ {
                if rf.log[i].Term == reply.ConflictTerm {
                    reply.ConflictIndex = i
                    break
                }
            }
		}		
        return
	}
}
~~~





















