## РОСДОМОФОН - Открытие по sms

### Нужно:
- 3g/4g или какой-нить еще модем
- Сервачок

### Установка:
```bash
wget https://github.com/mimimix/rosdomofon-sms/raw/refs/heads/main/docker-compose.yml
wget https://github.com/mimimix/rosdomofon-sms/raw/refs/heads/main/conf.empty.yml
mv conf.empty.yml conf.yml
# Заполняем конфиг
sudo docker-compose up -d
```

### Конфиг:
```
SECRET_KEY - рандомные символы
PROTECTION_CODE - слово, которое обязательно должно быть в смс (что-то вроде пароля)
KEY_ID - перехватываем http запрос приложения к https://rdba.rosdomofon.com/rdas-service/api/v1/temporary_keys и берем из тела запроса
HTTP_PORT - внутренний порт контейнера, ни на что не влияет
MODEM_URL - http путь до модема
LAST_SMS_FILE - файл с последними номерами смс (название, ни на что не влияет)
SMS_ALIVE_TIME - если смс отправлено ранее, чем указанное кол-во секунд - скипаем
REFRESH_TOKEN - перехватываем http запрос приложения к https://rdba.rosdomofon.com/authserver-service/oauth/token и берем из тела запроса
```

### Сложности:
#### 1. Модем  
Код написан под конкретно мой, так что вам скорее всего придется переписывать работу с модемом под свой.
#### 2. Настройка модема
В моем случае нужно было прокинуть запросы на адрес 192.168.8.1 в модем, но чтобы инет через него не шёл.
У меня решилось добавлением такого в rc.local:
```bash
ip addr add 192.168.8.2/32 dev enx01e101f00
ip link set enx01e101f00 up
ip route add 192.168.8.1/32 dev enx01e101f00
```

### Как работает:
Отправляетем смс в виде: "domofon PROTECTION_CODE" и дверь открывается

### Использованные библиотеки:
Библиотека для работы с модемом (переделал под себя): https://github.com/lagarciag/huaweimodem/tree/main