# Sample project
---

## Запуск

* Склонировать проект и перейти в папку проекта

	`git clone git@github.com:kshamiev/sample.git && cd sample`

* Выполнить обновление зависимостей командой:

	`make dep`
	
* Перейти в папку `conf` и скопировать файлы:

	`sample.example.log.yml -> log.yml`
		
	`sample.example.yml -> sample.yml`
	
* Отредактировать файлы `log.yml` и `sample.yml`

* Создать пустую базу данных и доступ к ней

* Прописать доступ к базе данных в `sample.yml`

* Установить утилиту `https://gopkg.in/webnice/migrate.v1` на сервер в папку доступную в `${PATH}`

* Скомпилировать проект командой:

	`make build`
	
* Запуск проекта в режиме продакшн (нет вывода логов в STDERR/STDOUT):

	`make run`
	
* Запуск проекта в development режиме:

	`make dev`
