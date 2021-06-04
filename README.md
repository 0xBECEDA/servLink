Реализация сервиса коротких ссылок.

Файл содержат:
- servShortLink.go - бэкэнд сервиса
- servShortLink_test.go - тесты сервиса
- front.html - фронтенд сервиса
- doc.org - документация

Короткий гайд:
- скачайте все файлы
- скомпилируйте файл servShortLink.go, использовав команду make
- запустите бэкэнд сервиса, введя в консоль ./servShortLink
- откройте браузер и перейдите на http://localhost:8080/serv_link -
  вам откроется форма для ввода короткой ссылки
- введите ваш тестовый адрес для получения короткой ссылки в первую
  форму (например, https://www.google.com ) нажмите кнопку submit
- перейдите по вернувшейся вам короткой ссылке - у вас должна
  открыться заданная вами страница
- для просмотра статистики введите в адресную строку браузера
  http://localhost:8080/get_link_statistiсs/?get_link_statistiсs= + вашу короткую ссылку
  Например,
  http://localhost:8080/get_link_statistiсs/?http://localhost:8080/?u=rGu2
- вы должны получить ответ от сервиса в консоли, сколько раз вы совершили
  переход по заданной короткой ссылке

Для запуска тестов используйте команду go test -v