#### What is got
- 'got' is a command line keeper of your notes, jobs, deadlines, lists, and time tracking.
- 'got' is a prototype, written in swift and bash scripts
- 'got' wants to know what you think of it! email any feedback to vic@vixac.com


#### Examples

- **got to Try out Got**       _(adds the item 'Try out Got'. No list, due for today.)_
- **got it shop Strawberry icecream**      _(adds the item 'Strawberry icecream' to the list 'shop'. No deadline.)_
- **got till 19th holiday book flights**      _(adds 'book flights' to the list 'holiday' with deadline 19th of this or next month, whichever is coming up next)_
- **got till 21/05/18 admin Email Bob about the thing**        _(adds 'Email Bob about the thing' to the list 'admin' for the 21st May 2018)_
- **got jobs admin**      -(show all jobs under the admin list)


#### Getting Started

You'll need the swift compiler (swiftc) which can be found at https://swift.org/download/#releases
If you're running on Windows, you can get started [here.](https://www.youtube.com/watch?v=dQw4w9WgXcQ)

**Installation**
- clone this repo


- for simple installation, run this:
```
 sudo ./simple_install.sh
```
Alternatively, set $GOT to be the directory where you want to keep your 'got' data, then:
- run install.sh  (creates the got executable in the directory of the repository)
- put the got executable in your :path, (you can  run sudo move.sh to put it in /usr/local/bin if you like).

#### How does it work
- 'got' stores everything in plain text files in your $GOT directory, so if you want your got synced across machines, all you need to do is put your $GOT in something like google drive or dropbox
- 'got' is mostly entirely keyboard based, except when you want to do something to an item, in which case you'll want to copy paste the hash (double click CTRL+C CTRL+V)
- 'got' uses colours, so works best with a dark or black background
- 'got', like 'echo', lets you enter plain text into your terminal. Also, like 'echo', if you're going to use reserved symbols like ' or .*, you'll want to use quotes for the description of your item. For example: 'got to test "using got"' and 'got till 3rd "Test every got command"'

#### Items, lists, hashes
Every job or task you put in 'got' is an item. It's your description, along with a hash, the time you created it, and the list and deadline if there is one.
A hash is just a unique identifier for the item, you can use it to do things to the item like mark it complete, or track time and keep notes on it.
A list can be any word as long as there are no underscores. You can group lists with the prefix, so if you had lists called 'some', 'something' and 'something-else', you could look up items in all 3 by typing 'got jobs some'.

### Usage
All the commands are described best in:
```
got help
```

### Work in progress
There are many parts still in need of more work:
- 'got complete' wants to show all work you've done recently
- 'got' wants to sync to an app on your phone as well. It's on the cards.
- 'got' wants to support Bloomberg style number<GO>s when selecting from a list
- 'got' wants to be easier to install and get started
- 'got' wants to be better at coping with merge conflicts over file sharing platforms
- 'got' wants to be awesome for note keeping
- 'got' wants to support multiple views, for multiple areas of life
- 'got' wants to reduce the need to copy paste hashes as much as possible.
