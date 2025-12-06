PING request:

```
*1\r\n$4\r\nPING\r\n
```

PONG response:

```
+PONG\r\n
```

SET request:

```
*3\r\n
$3\r\nSET\r\n
$5\r\nmykey\r\n
$5\r\nvalue\r\n
```

SET response:

```
+OK\r\n
```

GET request:

```
*2\r\n
$3\r\nGET\r\n
$5\r\nmykey\r\n
```

GET response when key exists:

```
$5\r\nvalue\r\n
```

GET response when key does not exist (RESP2):

```
$-1\r\n
```

DEL request:

```
*2\r\n
$3\r\nDEL\r\n
$5\r\nmykey\r\n
```

DEL response (integer count of removed keys):

```
:1\r\n
```

INCR request:

```
*2\r\n
$4\r\nINCR\r\n
$3\r\nctr\r\n
```

INCR response:

```
:1\r\n
```

EXISTS request:

```
*2\r\n
$6\r\nEXISTS\r\n
$5\r\nmykey\r\n
```

EXISTS response:

```
:1\r\n
```

TYPE request:

```
*2\r\n
$4\r\nTYPE\r\n
$5\r\nmykey\r\n
```

TYPE response:

```
+string\r\n
```

hello 2

```
*14\r\n
$6\r\n
server\r\n
$5\r\n
redis\r\n
$7\r\n
version\r\n
$5\r\n
6.3.4\r\n
$5\r\n
proto\r\n
:2\r\n
$2\r\n
id\r\n
:4\r\n
$4\r\n
mode\r\n
$9\r\n
standalone\r\n
$4\r\n
role\r\n
$6\r\n
master\r\n
$7\r\n
modules\r\n
*0\r\n````