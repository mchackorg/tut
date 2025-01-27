package main

var conftext = `# Configuration file for tut

[general]
# If the program should check for new toots without user interaction. If you
# don't enable this the program will only look for new toots when you reach the
# bottom or top of your feed. With this enabled it will check for new toots
# every x second.
# default=true
auto-load-newer=true

# How many seconds between each pulling of new toots if you have enabled
# auto-load-newer
# default=60
auto-load-seconds=60

# The date format to be used. See https://godoc.org/time#Time.Format
# default=2006-01-02 15:04
date-format=2006-01-02 15:04

# Format for dates the same day. See date-format for more info.
# default=15:04
date-today-format=15:04

# This displays relative dates instead for statuses that are one day or older
# the output is 1y2m1d (1 year 2 months and 1 day)
# 
# The value is an integear
# -1     = don't use relative dates
#  0     = always use relative dates, except for dates < 1 day
#  1 - ∞ = number of days to use relative dates
# 
# Example: date-relative=28 will display a relative date for toots that are
# between 1-28 days old. Otherwhise it will use the short or long format.
# default=-1
date-relative=-1

# The timeline that opens up when you start tut.
# Valid values: home, direct, local, federated
# default=home
timeline=home

# The max width of text before it wraps when displaying toots.
# 0 = no restriction.
# default=0
max-width=0

# If you want to display a list of notifications under your timeline feed.
# default=true
notification-feed=true

# Where do you want the list of toots to be placed?
# Valid values: left, right, top, bottom.
# default=left
list-placement=left

# If you have notification-feed set to true you can display it under the main
# list of toots (row) or place it to the right of the main list of toots
# (column).
# default=row
list-split=row

# Hide notification text above list in column split. It's displayed as
# [N]otifications.
# default=false
hide-notification-text=false

# You can change the proportions of the list view in relation to the content
# view list-proportion=1 and content-proportoin=3 will result in the content
# taking up 3 times more space.
# Must be n > 0
# default=1
list-proportion=1

# See list-proportion
# default=2
content-proportion=2

# If you always want to quote original message when replying.
# default=false
quote-reply=false

# If you're on an instance with a custom character limit you can set it here.
# default=500
char-limit=500

# If you want to show icons in the list of toots.
# default=true
show-icons=true

# If you've learnt all the shortcut keys you can remove the help text and only
# show the key in tui. So it gets less cluttered.
# default=false
short-hints=false

# If you want to show a message in the cmdbar on how to access the help text.
# default=true
show-help=true

# If you don't want the whole UI to update, and only the text content you can
# set this option to true. This will lead to some artifacts being left on the
# screen when emojis are present. But it will keep the UI from flashing on every
# single toot in some terminals.
# default=true
redraw-ui=true

[media]
# Your image viewer.
# default=xdg-open
image-viewer=xdg-open

# Open the image viewer in the same terminal as toot. Only for terminal based
# viewers.
# default=false
image-terminal=false

# If images should open one by one e.g. "imv image.png" multiple times. If set
# to false all images will open at the same time like this "imv image1.png
# image2.png image3.png". Not all image viewers support this, so try it first.
# default=true
image-single=true

# If you want to open the images in reverse order. In some image viewers this
# will display the images in the "right" order.
# default=false
image-reverse=false

# Your video viewer.
# default=xdg-open
video-viewer=xdg-open

# Open the video viewer in the same terminal as toot. Only for terminal based
# viewers.
# default=false
video-terminal=false

# If videos should open one by one. See image-single.
# default=true
video-single=true

# If you want your videos in reverse order. In some video apps this will play
# the files in the "right" order.
# default=false
video-reverse=false

# Your audio viewer.
# default=xdg-open
audio-viewer=xdg-open

# Open the audio viewer in the same terminal as toot. Only for terminal based
# viewers.
# default=false
audio-terminal=false

# If audio should open one by one. See image-single.
# default=true
audio-single=true

# If you want to play the audio files in reverse order. In some audio apps this
# will play the files in the "right" order.
# default=false
audio-reverse=false

# Your web browser.
# default=xdg-open
link-viewer=xdg-open

# Open the browser in the same terminal as toot. Only for terminal based
# browsers.
# default=false
link-terminal=false

[open-custom]
# This sections allows you to set up to five custom programs to upen URLs with.
# If the url points to an image, you can set c1-name to img and c1-use to imv.
# If the program runs in a terminal and you want to run it in the same terminal
# as tut. Set cX-terminal to true. The name will show up in the UI, so keep it
# short so all five fits.
# 
# c1-name=img
# c1-use=imv
# c1-terminal=false
#    
# c2-name=
# c2-use=
# c2-terminal=false
#   
# c3-name=
# c3-use=
# c3-terminal=false
#   
# c4-name=
# c4-use=
# c4-terminal=false
#   
# c5-name=
# c5-use=
# c5-terminal=false

[open-pattern]
# Here you can set your own glob patterns for opening matching URLs in the
# program you want them to open up in. You could for example open Youtube videos
# in your video player instead of your default browser.
# 
# You must name the keys foo-pattern, foo-use and foo-terminal, where use is the
# program that will open up the URL. To see the syntax for glob pattern you can
# follow this URL https://github.com/gobwas/glob#syntax. foo-terminal is if the
# program runs in the terminal and should open in the same terminal as tut
# itself.
# 
# Example for youtube.com and youtu.be to open up in mpv instead of the browser.
# 
# y1-pattern=*youtube.com/watch*
# y1-use=mpv
# y1-terminal=false
# 
# y2-pattern=*youtu.be/*
# y2-use=mpv
# y2-terminal=false

[desktop-notification]
# Notification when someone follows you.
# default=false
followers=false

# Notification when someone favorites one of your toots.
# default=false
favorite=false

# Notification when someone mentions you.
# default=false
mention=false

# Notification when someone boosts one of your toots.
# default=false
boost=false

# Notification of poll results.
# default=false
poll=false

# Notification when there is new posts in current timeline.
# default=false
posts=false

[style]
# All styles can be represented in their HEX value like #ffffff or with their
# name, so in this case white. The only special value is "default" which equals
# to transparent, so it will be the same color as your terminal.
# 
# You can also use xrdb colors like this xrdb:color1 The program will use colors
# prefixed with an * first then look for URxvt or XTerm if it can't find any
# color prefixed with an asterik. If you don't want tut to guess the prefix you
# can set the prefix yourself. If the xrdb color can't be found a preset color
# will be used.

# The xrdb prefix used for colors in .Xresources.
# default=guess
xrdb-prefix=guess

# The background color used on most elements.
# default=xrdb:background
background=xrdb:background

# The text color used on most of the text.
# default=xrdb:foreground
text=xrdb:foreground

# The color to display sublte elements or subtle text. Like lines and help text.
# default=xrdb:color14
subtle=xrdb:color14

# The color for errors or warnings
# default=xrdb:color1
warning-text=xrdb:color1

# This color is used to display username.
# default=xrdb:color5
text-special-one=xrdb:color5

# This color is used to display username and key hints.
# default=xrdb:color2
text-special-two=xrdb:color2

# The color of the bar at the top
# default=xrdb:color5
top-bar-background=xrdb:color5

# The color of the text in the bar at the top.
# default=xrdb:background
top-bar-text=xrdb:background

# The color of the bar at the bottom
# default=xrdb:color5
status-bar-background=xrdb:color5

# The color of the text in the bar at the bottom.
# default=xrdb:foreground
status-bar-text=xrdb:foreground

# The color of the bar at the bottom in view mode.
# default=xrdb:color4
status-bar-view-background=xrdb:color4

# The color of the text in the bar at the bottom in view mode.
# default=xrdb:foreground
status-bar-view-text=xrdb:foreground

# Background of selected list items.
# default=xrdb:color5
list-selected-background=xrdb:color5

# The text color of selected list items.
# default=xrdb:background
list-selected-text=xrdb:background

`
