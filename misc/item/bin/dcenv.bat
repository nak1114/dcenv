@echo off

if not DEFINED DCENV_COMMAND (
    set DCENV_COMMAND=dcenv
)

if not DEFINED DCENV_SHELL (
    set DCENV_SHELL=windows
)

IF EXIST "%DCENV_HOME%\tmp\%DCENV_COMMAND%.bat" (
  del /F /Q "%DCENV_HOME%\tmp\%DCENV_COMMAND%.bat"
)
call %~dp0..\files\dcenv.exe %*

set DCENV_ARGS=
IF EXIST "%DCENV_HOME%\files\__args__" (
  set /p DCENV_ARGS=<"%DCENV_HOME%\files\__args__"
)
IF EXIST "%DCENV_HOME%\tmp\%DCENV_COMMAND%.bat" (
  "%DCENV_HOME%\tmp\%DCENV_COMMAND%.bat" %DCENV_ARGS%
)

