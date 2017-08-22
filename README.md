
### Installation

clone this repo

Choose where you want to keep your Got contents files, then add these 3 lines to your .profile

```
export GOT_SRC=<path/to/this/repo>
export GOT_CONTENTS=<path/to/wherever/you/want/your/contents>
. $GOT_SRC/got_env
```
Sourcing this file will also compile Got if the executable does not exist yet.
You should be set! 


### Usage
```
 got help
```
Command   Parameters                            Summary                                 More
it        <item>                                quickly keep an item                    The simplest thing you can add: create an item with no list or deadline. For example: 'got it Try out Got'
to        <list> <item>                         add an item to a list                   Add an item to a list. For example: 'got to shop Strawberry icecream' creates a list called 'shop' and adds an item to it called: 'Strawberry Icecream'
till      [offset|nth|dd/mm/yy] <list> <item>   add an item with a deadline to a list   Add an item with a deadline. There are 3 ways to do this. 'got till 5 read Grapes of Wrath' creates a deadline in 5 days time. 'got till 19th holiday book flights' creates a deadline
                                                                                        either the 19th of this month, or next month if 19th has passed. 'got till 21/05/18 admin Email Bob about the thing' sets the deadline to 21st May 2018

jobs      [list]                                see everything                          Shows all active jobs. You can specify a list. For example: 'got jobs accounts'
what                                            see summaries of lists
start     <hash>                                time an item, and take notes            Starts the timer on an item. For example: got start 0c441b2b0 It will block this shell. You can add notes to this hash by writing into the shell with the timer. When you are finished, CTRL C or type 'stop'
done      <hash>                                mark an item complete                   got done 0c441b2b0
remove    <hash>                                erase an item                           got remove 0c441b2b0
today     <list> <item>                         add an item to do today



### Work in progress
There are many parts still in need of more work.
