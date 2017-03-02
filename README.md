# dcenv : Help your develop environment with docker.

Dcenv supports installation and execution of commands by docker.
Also, docker makes it easy to change the version of the command.
By using the docker container,
package management independent of OS is possible, unlike apt-get and yum.


## How It Works

The dcenv command consulted rbenv.
In the example we install the `go` command.

#### 1.Install the configuration.

    ~~~ 
    $ dcenv install golang
    ~~~

#### 2.Deploy the configuration.

    ~~~ 
    $ dcenv local golang
    $ dcenv rehash
    ~~~


#### 3.Execute the command

    ~~~ 
    $ go version
    go version go1.8 linux/amd64
    ~~~

Only this.

### Switching environment

#### Switching versions

    ~~~ 
    $ dcenv tag golang:1.7
    $ go version
    go version go1.7.5 linux/amd64
    ~~~

You could change to version 1.7.

#### Switching dir

The command can be executed under the directory where `local` is executed.

    ~~~ 
    $ mkdir child
    $ cd child
    $ go version
    go version go1.7.5 linux/amd64
    
    $ cd ../../
    $ go version
    bash: go: command not found
    ~~~

### More dcenv commands

There are several commands for checking the dcenv environment.



#### Retrieve the command list registered in the registry

    ~~~ 
    $ dcenv install --list
    golang           
    naktak/clang     
    naktak/dcenv_test
    
    ~~~

#### Retrieve installed command list

    ~~~ 
    $ dcenv yard --list
    golang           
    
    ~~~

#### Acquire a list of executable commands in the current directory

    ~~~ 
    $ dcenv command --list
    COMMAND   | IMAGE        (FILE)
    go        | golang       (/tmp/hoge/.dcenv_bash)
    
    ~~~

#### Show details of installed command

    ~~~ 
    $ dcenv yard golang
    ---[ 0 ]---------                                                                                     
      Id      :                                                                                           
      Owner   :                                                                                           
      Image   : golang                                                                                    
      Brief   : golang for go command                                                                     
      Desc    :                                                                                           
      Pri     : 0                                                                                         
      Config  : bash                                                                                      
        Tag     :                                                                                         
        Fake    : false                                                                                   
        Commands: map[go:map[]]                                                                           
        Script  :                                                                                         
    docker run --rm -it -v "$(pwd)":/myapp -w /myapp golang:{{.Tag}} go "$@"                              
                                                                                                          
      Config  : windows                                                                                   
        Tag     :                                                                                         
        Fake    : false                                                                                   
        Commands: map[go:map[]]                                                                           
        Script  :                                                                                         
    docker run --rm -it -v %CD%:/myapp -w /myapp -e GOOS=windows -e GOARCH=386 golang:{{.Tag}} go %*      
    
    ~~~

#### Display details of executable commands in the current directory

    ~~~ 
    $ dcenv command go
    ---[ go ]------------------------------
    (Short)
    
    ~~~ 


## How to create a command

There are still few commands registered in the registry.
Please also create your command and register it in the registry.

#### Create a command template.

    ~~~ 
    $ dcenv yard --create busybox
    Edit: ~/.dcenv/image_yard/busybox.yml
    
    ~~~ 
#### Edit this template.

    ~~~ 
    $ vi ~/.dcenv/image_yard/busybox.yml
    
    ~~~ 

[For details of editing refer to here.](./docs/script.md)
or `dcenv install naktak/dcenv-script-sample`

#### Deploy your command and check.

    ~~~ 
    $ dcenv local busybox
    $ dcenv rehash
    $ ls

#### Sign up for the registry.

Click here to sign up.
https://nak1114.github.io/dcenv/sign_up.html


#### Register your command in the registry.

    ~~~ 
    $ dcenv push busybox
    
    ~~~ 


#### Log out from the registry

    ~~~ 
    $ dcenv logout
    
    ~~~ 

#### Delete your\command from the registry

    ~~~ 
    $ dcenv push -d busybox
    
    ~~~ 

## DCEnv Cheat Sheet

![DCEnv cheat sheet](./docs/DCEnv_CheatSheet.png "DCEnv cheet sheat")

## How to Install dcenv

### Linux(Bash)

Specify the environment variable in .bashrc. And reload this file.

    ~~~ .bashrc
    export DCENV_HOME=~/.dcenv
    export PATH=${DCENV_HOME}/shims;${DCENV_HOME}/bin;${PATH}
    ~~~ 

Download the file and extract it.

    ~~~ 
    $ wget --no-check-certificate https://github.com/nak1114/dcenv/releases/download/v0.0.1/dcenv-v0.0.1-linux-amd64.tar.gz

    $ mkdir -p $DCENV_HOME
    $ tar xvfz dcenv-0.0.1-linux-amd64.tar.gz -C $DCENV_HOME
    ~~~ 


### Windows

Specify the environment variable in system. And restart console.

    ~~~ 
    > set DCENV_HOME=%USERPROFILE%\.dcenv
    > set PATH=%DCENV_HOME%\shims;%DCENV_HOME%\bin;%PATH%
    ~~~ 

Download the file and extract it.

    ~~~ 
    > powershell wget https://github.com/nak1114/dcenv/releases/download/v0.0.1/dcenv-v0.0.1-windows-amd64.zip
    > md %DCENV_HOME%
    > powershell Expand-Archive dcenv-0.0.1-windows-amd64.zip  -DestinationPath %DCENV_HOME%
    ~~~ 


## Environment variables

You can affect how rbenv operates with the following settings:

name | default | description
-----|---------|------------
`DCENV_HOME` | `dirname $0`/.. | Defines the directory under which DCEnv commands and shims reside.
`DCENV_DIR` | `$PWD` | Directory to start searching for `.dcenv_*` files.
`DCENV_SHELL` | `bash` or<BR/> `windows` | Defines your command shell.
`DCENV_COMMAND` | | Internal use only
`DCENV_ARGS` | | Internal use only


## Contribution

1. Fork it ( http://github.com/nak1114/dcenv/fork )
2. Create your feature branch (git checkout -b my-new-feature)
3. Commit your changes (git commit -am 'Add some feature')
4. Push to the branch (git push origin my-new-feature)
5. Create new Pull Request

## Licence

MIT
