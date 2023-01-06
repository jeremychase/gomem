# playing with pprof and memory 'leaks'

# usage

To use pprof we need to build and use the binary:

```
$ go build .
$ ./gomem -iterations 1000000
writing output-mid.mprof at 500000
writing output-end.mprof at 1000000
```

then if you want to see allocations:

```
$ go tool pprof -sample_index=alloc_space gomem output-mid.mprof 
File: gomem
Type: alloc_space
Time: Jan 6, 2023 at 11:19am (EST)
Entering interactive mode (type "help" for commands, "o" for options)
(pprof) top 2
Showing nodes accounting for 112.41MB, 56.51% of 198.91MB total
Dropped 12 nodes (cum <= 0.99MB)
Showing top 2 nodes out of 18
      flat  flat%   sum%        cum   cum%
   58.01MB 29.16% 29.16%   112.01MB 56.31%  encoding/json.Unmarshal
   54.40MB 27.35% 56.51%   196.91MB 98.99%  main.leak
(pprof) list main.leak
Total: 198.91MB
ROUTINE ======================== main.leak in /Users/jchase/src/github.com/jeremychase/gomem/main.go
   54.40MB   196.91MB (flat, cum) 98.99% of Total
         .          .     60:   }
         .          .     61:
         .          .     62:   t := data{rand.Int(), rand.Int()}
         .          .     63:
         .          .     64:   // Marshal/Unmarshal into new value for spurious allocations
    5.50MB     5.50MB     65:   var r data
         .          .     66:   {
       5MB    35.50MB     67:           js, err := json.Marshal(t)
         .          .     68:           if err != nil {
         .          .     69:                   log.Println("marshal failure")
         .          .     70:                   return err
         .          .     71:           }
         .          .     72:
         .   112.01MB     73:           if err := json.Unmarshal(js, &r); err != nil {
         .          .     74:                   log.Println("unmarshal failure")
         .          .     75:                   return err
         .          .     76:           }
         .          .     77:   }
         .          .     78:
         .          .     79:   // 'leak'
   43.90MB    43.90MB     80:   storage[r.Key] = r.Value
         .          .     81:
         .          .     82:   return nil
         .          .     83:}
         .          .     84:
         .          .     85:func parseFlags() (*os.File, *os.File) {
(pprof) 
```

or memory in use:


```
go tool pprof -sample_index=inuse_space gomem output-mid.mprof 
File: gomem
Type: inuse_space
Time: Jan 6, 2023 at 11:19am (EST)
Entering interactive mode (type "help" for commands, "o" for options)
(pprof) top2
Showing nodes accounting for 32.69MB, 98.49% of 33.19MB total
Showing top 2 nodes out of 21
      flat  flat%   sum%        cum   cum%
   31.19MB 93.97% 93.97%    31.19MB 93.97%  main.leak
    1.50MB  4.52% 98.49%     1.50MB  4.52%  runtime.allocm
(pprof) list main.leak
Total: 33.19MB
ROUTINE ======================== main.leak in /Users/jchase/src/github.com/jeremychase/gomem/main.go
   31.19MB    31.19MB (flat, cum) 93.97% of Total
         .          .     75:                   return err
         .          .     76:           }
         .          .     77:   }
         .          .     78:
         .          .     79:   // 'leak'
   31.19MB    31.19MB     80:   storage[r.Key] = r.Value
         .          .     81:
         .          .     82:   return nil
         .          .     83:}
         .          .     84:
         .          .     85:func parseFlags() (*os.File, *os.File) {
(pprof) 
```