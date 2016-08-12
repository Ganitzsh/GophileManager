# GophileManager

GophileManager is a web app allowing you easily manipulate your files in a specific directory.

The motivation behind this project is simple, I wanted to use `Go` - my favorite language - to build something that I would use daily. I'm very often manipulating files from my seedbox.

I was using nginx to access a specific directory and get some files from it. It was not convenient and very ugly. Plus, if I wanted to do a little bit of file management I had to connect to my server through SSH.

And then, one day, I realized it was 2016 and that my knowledge of `Go` was more than enough to do something cool and clean to avoid spending too much time doing annoying things.

So I developed **GophileManager**!

## Note

This is not a finished product. I'm only working on it in my free time. I am not responsible for the bad things that might happen to you while using it. Use at your own risks.

## Installation

First, you have to clone the repository **inside your GOPATH**. Then, you can either directly run it using the revel CLI:

```shell
revel run github.com/ganitzsh/WebManager <run_mod> # prod or dev
```

Or you can build it with:

```shell
# target_dir is the directory where to put the built files
revel build github.com/ganitzsh/WebManager <target_dir> <run_mode>
```

Then, you'll have to run the app using the `run.sh` script inside the freshly created `<target_dir>` and you should be ready to configure!

### Configuration

#### Configuration file

In order to fully experience the greatness of **GophileManager**, the `app.conf` file needs to be updated.

If you run the app directly from the cloned repository, you'll find this file in `<repo>/confing/app.conf`.

If you chose to build the app, you can find it in `<target_dir>/src/github.com/ganitzsh/WebManager/conf`

Open it with your beloved text editor and find the following part:

```
# app.main_dir = /tmp
# app.trash_dir = /tmp/trash
# app.host = http://localhost
```

- `app.main_dir` is the directory you want to manage
- `app.trash_dir` is the directory that's ought to be used a trash (not kewl!)

  - **NOTE:** If not set, the trash system won't be available and you won't be able to delete the files.

- `app.host` is VERY IMPORTANT. This is the URL from which the clients are going to call the app. E.g: you want to type this address in your browser: `http://manage.domain.com` then you will have to write: `app.host = http://manage.domain.com`. It also works with `https://`.

  - **NOTE:** The reason for this is the CORS specification, if this variable isn't set correctly the websockets won't work and no feedback will be given to you.

These 3 variables have to be set to fully enjoy **GophileManager**.

#### Script file

The client-side of the app needs to connect to the Socket.IO server.

To do so, you must edit the `script.js` file located here: `<repo>/public/js/script.js` and find the following line:

```javascript
socket = io('http://localhost:9000');
```

Replace `'localhost'` with your server's host. If your domain name is `www.domain.com` then replace write

```javascript
socket = io('http://www.domain.com:9000');
```

If you changed the port number inside the `app.conf` file, make sure to update the port as well.

**NOTE:** It doesn't seem to work behind a proxy. I tested it with an nginx proxy redirecting the calls to the port `9000` to a virtualhost, but no luck so far.

## All set!

The default port of the app is `9000`. Once it is up and running You can then access **GophileManager** here:

```
http://localhost:9000/app
```

**NOTE:** The `/app` needs to be added!

## What is it? And what can I achieve wit this?

**GophileManager** is primarily a file manager.

The main framework used to run the app is [Revel](https://revel.github.io/), coupled with the `Go`'s [Socket.IO](https://github.com/googollee/go-socket.io) implementation.

### Features

With **GophileManager** you can achieve few operations on files:

- Delete
- Compress (as a `.tar` ball)
- Download
- Navigate through the directories

#### UI

**GophileManager** comes with a nice responsive GUI made with bootstrap.

#### File analysis

The files contained in the directory managed by **GophileManager** and all it's children are analyzed in order to dynamically separate the files in `categories`.

Each `category` has `sub-categories` except for the directories (yes, they are files as well! But special ones)

#### Trash system

**GophileManager** comes with a trash system. It is a directory that's going to be used as a trash. The unneeded files are placed into the trash before you can erased them from the disk.

This feature can be disabled, but the main purpose is to avoid mistakes! Once a file is deleted, it's gone!

The trash can be emptied with a single push of a button.

#### Notifications

The server is using `SocketIO` to give feedback to the user for time-consuming operations. The page can be reloaded during an operation (e.g: compression, conversion), and a notification will be sent to the clients when it's done or in case of failure.

#### Navigation

As the main purpose of **GophileManager** is to manage a directory, you can easily navigate through it's content. Some navigation options are easily reachable such as:

- Reload the content
- Going back to parent (..)
- Back in the main directory (home)
- Access the trash
- Empty the trash

#### Common operations

Some common file operations from **GophileManager**:

- Move to trash
- Compress
- Download

#### Compression

With **GophileManager**, you can easily create an archive from any file.

Simply press the button, a dialog is going to ask you for the new name and create a `.tar` archive right away!

#### Special operations

Some special operations are made possible for particular file formats.

##### Images

**GophileManager** Offers the possibility to preview an image.

##### Videos

A few formats of video are recognized by **GophileManager**.

If the video is in the `mp4` format, you can play it directly from your browser! Awesome right?

If the `ffmpeg` program is installed on your machine, **GophileManager** allows you to convert any video into the `mp4` format. The new file will be placed in the same directory with the same name but with a `.mp4` extension.

**NOTE:** The video won't be re-encoded. `ffmpeg` will use the original codec but will rewrite the video as an `mp4` file. The main reason is a matter of performance of the hosting machine, re-encoding a video takes a lot of time and CPU resources. It will work most of the time, but I had a lot of troubles with the `avi` format. Just try and see what happens.

## Developer

Each feature is accessible from a specific route, you can then develop your own client.

**NOTE:** More documentation will come in the future. For now, take a look at the code!
