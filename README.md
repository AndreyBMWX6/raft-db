# raft-db 
[![Build Status](https://travis-ci.org/AndreyBMWX6/raft-db.svg?branch=master)](https://travis-ci.org/AndreyBMWX6/raft-db)
##Информация для защиты:
Презентация и видео, демонстрирующие работу программы находятся в папке `materials`

## Руководство
## Как запустить кластер:
Перейдите в папку `internal/config` и установите диапазон используемых портов в файле `config.go`  
Пример:  
```
// *SET PORTS RANGES HERE*

	// ip ports range
	firstIPPort := "8001"
	lastIPPort  := "8006"

	// URL ports range
	firstURLPort := "8081"
	lastURLPort := "8086"
```

### Для пользователей Windows:
После этого запустите командную строку,  
в ней, используя команду `cd <Path>` перейдите в папку `raft-db/scripts` 
и запустите скрипт `run-all.bat` следующим образом:  
```
Win + R
cmd
Enter

cd <path to raft-db>/scripts
run-all.bat <ip> <first ip port> <last ip port> <first url port> <last url port>

Например:
run-all.bat 127.0.0.1 8001 8006 8081 8086
```

#### Примечание:
Количество url и ip портов должно совпадать, диапазоны портов должны совпадать с диапазонами в `config.go`,
порт являющийся началом диапазона должен быть меньше порта являющегося концом диапазона, 
также диапазоны не должны пересекаться. 
Ниже приведены примеры некорректного использования скрипта: 
```
// Диапазон выделенный под ip порты больше диапазона под url порты
run-all.bat 127.0.0.1 8001 8006 8081 8083
// Недостаточное количество параметров
run-all.bat 127.0.0.1 8001 8006 8081
run-all.bat 127.0.0.1 8001 8006
run-all.bat 127.0.0.1
run-all.bat
```

### Как запустить 1 сервер:
Если вы не использовали все порты из диапазона портов в файле `config.go`,
то вы можете запустить эти сервера по одному используя `run.bat`.
```
cd <path to raft-db>/scripts
start run.bat <ip> <ip port> <url port> <run router> <username>

Примеры:
start run.bat 127.0.0.1 8001 8081 false user1
start run.bat 127.0.0.1 8002 8082 true user2
```

#### Примечание:
Во избежание ошибок рекомендуется устанавливать значение параметра `<run router>` равным `'false'`.
Вы можете `задать значение параметра равным 'true' только если вы не собираетесь запускать скрипт run-all.bat!`
А вместо этого хотите запустить все сервера кластера по одному. 
Тогда следует `установить параметр равным 'true' только для 1 сервера!`  
Например в конфиг-файле есть следующие диапазоны:
```
// *SET PORTS RANGES HERE*

	// ip ports range
	firstIPPort := "8001"
	lastIPPort  := "8004"

	// URL ports range
	firstURLPort := "8081"
	lastURLPort := "8084"
```
Тогда мы можем сначала запустить 3 сервера, а после этого ещё 1:
```
// OK
run-all.bat 127.0.0.1 8001 8003 8081 8083
start run.bat 127.0.0.1 8004 8084 false user4
// Ошибка
run-all.bat 127.0.0.1 8001 8003 8081 8083
start run.bat 127.0.0.1 8004 8084 true user4 
// Почему: мы уже запустили роутер в run-all.bat
```

### Как остановить кластер:
Для остановки работы серверов перейдите в папку `raft-db/scripts`  
и запустите скрипт `run-all.bat`:
```
cd <path to raft-db>/scripts
kill-all.bat
```  

## Instructions
## How to run cluster:
Go to `internal/config` folder and set used ports range in `config.go` file.  
Example:  
```
// *SET PORTS RANGES HERE*

	// ip ports range
	firstIPPort := "8001"
	lastIPPort  := "8006"

	// URL ports range
	firstURLPort := "8081"
	lastURLPort := "8086"
```

### For Windows Users:
Run command line, and go to `raft-db/scripts` folder using `cd <Path>` command.
Run scripts `run-all.bat` this way:  
```
Win + R
cmd
Enter

cd <path to raft-db>/scripts
run-all.bat <ip> <first ip port> <last ip port> <first url port> <last url port>

Example:
run-all.bat 127.0.0.1 8001 8006 8081 8086
```

#### Note:
Amount of url и ip ports must fit each other, port ranges must fit with port ranges in `config.go`,
begin of range port must be less than end of range port, 
Also ranges mustn't cross. 
Here are examples of incorrect script using: 
```
// ip ports range is bigger than url ports range
run-all.bat 127.0.0.1 8001 8006 8081 8083
// not enough parametres were provided
run-all.bat 127.0.0.1 8001 8006 8081
run-all.bat 127.0.0.1 8001 8006
run-all.bat 127.0.0.1
run-all.bat
```

### How to run 1 server:
You can run some servers one by one if you didn't use all config range using `run.bat`.
```
cd <path to raft-db>/scripts
start run.bat <ip> <ip port> <url port> <run router> <username>

Examples:
start run.bat 127.0.0.1 8001 8081 false user1
start run.bat 127.0.0.1 8002 8082 true user2
```

#### Note:
It's recommended to set `<run router>` parameter to `'false'` to avoid errors.
You can `set this parameter to 'true' only if you are not going to execute run-all.bat script!`
and you are going to run all cluster servers one by one. 
Then you should `set this parameter to 'true' only for 1 server!`  
For example we have this config ranges:
```
// *SET PORTS RANGES HERE*

	// ip ports range
	firstIPPort := "8001"
	lastIPPort  := "8004"

	// URL ports range
	firstURLPort := "8081"
	lastURLPort := "8084"
```
So we can firstly run 3 servers and run one more after:
```
// OK
run-all.bat 127.0.0.1 8001 8003 8081 8083
start run.bat 127.0.0.1 8004 8084 false user4
// Error
run-all.bat 127.0.0.1 8001 8003 8081 8083
start run.bat 127.0.0.1 8004 8084 true user4 
// Why: we have already runned router in run-all.bat
```

### How to stop cluster:
Go to `raft-db/scripts` folder and run `run-all.bat`:
```
cd <path to raft-db>/scripts
kill-all.bat
```  

