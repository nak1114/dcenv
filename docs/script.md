# How to command writing


~~~
- id: ""
  owner: ""
  image: busybox
  brief: Busybox for ls and cat command
  desc: ""
  pri: 0
  config:
    bash:
      script: |
        tag: ""
        fake: false
        docker run --rm -it -v "$(pwd)":/myapp -w /myapp busybox:{{.Tag}} {{.Cmd}} "$@"
      commands: 
        ls: {}
        cat: {}
    windows:
      script: |
        docker run --rm -it -v "%CD%":/myapp -w /myapp busybox:{{.Tag}} {{.Cmd}} %*
      commands:
        ls: {}
        cat: {}
~~~

1. **id:**
This field is initially empty.
When you push, the registry will give.

2. **owner:**
This field is initially empty.
When you push, the registry will give.

3. **image:**
Docker containar image name.

4. **bash:**
Script for bash.

5. **windows:**
Script for windows.

6. **tag:**
container tag. you can change this field easy by `dcenv tag` command.default is `latest`.

7. **fake:**
Set to true if container image does not exist in DockerHub.
For example `naktak/dcenv-script-sample` should be true.

8. **script:**
This script is executed.

9. **commands:**
  command name with Env.

