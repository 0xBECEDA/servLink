package main

import
(
    "crypto/sha256"
    "fmt"
    "net/http"
    "encoding/base64"
    "net/url"
    "os/exec"
    "log"
)

// длина хэша
const hashLength = 4
// порт
const port = ":8080"
// базовый урл - нужен для построения коротких ссылок
var baseUrl = "http://localhost" + port + "/?u="

// структура - содердимое записи хэш-таблицы
// харнит в себе блинный урл и счетчик посещений
type fullUrl struct {
    Cnt int
    Url string
}

// таблица ссылок
var tableOfLinks = make( map[string]fullUrl )

// возвращает длинную ссылку для заданной короткой ссылки,
// если такая есть в системе.
// возвращает короткую ссылку и статус
func getFullUrl( shortUrl string ) (string, int) {

    val, ok := tableOfLinks[shortUrl]

    if ok == true {
        // увеличивает счетик посещений
        val.Cnt += 1
        tableOfLinks[shortUrl] = val
        return val.Url, http.StatusSeeOther
    }
    return "", http.StatusNotFound
}

// получает счетчик посещений короткой ссылки
func getLinkCnt( shortUrl string ) int {

    val, ok := tableOfLinks[shortUrl]
    if ok == true {
        return val.Cnt
    }
    return -1
}

// генерирует хэшот длинной ссылки и обрезает его до
// значения hashLength
func generateShortURL (link string) string {

    sum := sha256.Sum256( []byte( link ) )
    slice := sum[:]
    encoded := base64.StdEncoding.EncodeToString(slice)

    return encoded[:hashLength]
}

func getStat( w http.ResponseWriter, r *http.Request ) {
    hash, _ := getHashFromRequest( r )
    urlRequest := "http://localhost:8080/get_link_statistiсs?" + baseUrl + hash
    exec.Command("xdg-open", urlRequest).Start()
}

// получает статистику посещений для заданной в запросе короткой ссылке
// и отвечает на запрос либо строкой с кол-во посещений, либо ошибкой 404
func getLinkStatistiсs( w http.ResponseWriter, r *http.Request ) {

    hash, err := getHashFromRequest( r )

    if err == nil {
        encHash, _ := url.QueryUnescape( hash )
        cnt := getLinkCnt( encHash )

        if cnt >= 0 {
            str := fmt.Sprintf("Адрес %s посещали %d раз \n", baseUrl + hash, cnt)
            w.Write([]byte(str))

        } else {
            http.Redirect(w, r, baseUrl, http.StatusNotFound)
        }
    } else {
        http.Redirect(w, r, baseUrl, http.StatusNotFound)
        fmt.Println(err)
    }
}

func sentFront( w http.ResponseWriter, r *http.Request ) {
    http.ServeFile(w , r, "front.html")
}

func checkQuery( w http.ResponseWriter, r *http.Request ) {
    if len(r.URL.RawQuery) <= 0 {
        sentFront(w, r)
    } else {
        redirect(w, r)
    }
}

// завоит новую запись в хэш-таблице для заданной длинной ссылки
// и генерирует короткую ссылку
// вовзвращает короткую ссылку пользователю
func registerNewLink( w http.ResponseWriter, r *http.Request ) {

    query, _ := url.QueryUnescape(r.URL.RawQuery)

    if len(query) > 4 {
        url := query[4:]
        newStruct := fullUrl{ Url: url }
        hash := generateShortURL(url)
        tableOfLinks[hash] = newStruct

        str := fmt.Sprintf("Короткая ссылка для %s - %s \n", url, baseUrl + hash)
        w.Write([]byte(str))

    } else {
        w.Write([]byte("Введенный url неполный! \n"))
    }
}


// реализация ошибки получения хэша из параметров (query) запроса
type getHashFromRequestError struct {
    s string
}

func (e *getHashFromRequestError) Error() string {
    return e.s
}

// получает хэш из запроса
func getHashFromRequest( r *http.Request ) (string, error) {

    if (len(r.URL.RawQuery) > hashLength) {
        return r.URL.RawQuery[len(r.URL.RawQuery)- hashLength:], nil

    } else {
        err := getHashFromRequestError{s: "Невозможно получить короткую ссылку из запроса \n"}
        return "", &err
    }
}

// выполняет перенаправление с короткой ссылки на длинную
func redirect( w http.ResponseWriter, r *http.Request ) {

    link, err := getHashFromRequest( r )
    if err == nil {
        shortUrl, _ := url.QueryUnescape(link)
        url, statusCode := getFullUrl( shortUrl )

        if statusCode == 303 {
            http.Redirect(w, r, url, http.StatusSeeOther)

        } else {
            http.Redirect(w, r, shortUrl, http.StatusNotFound)
        }
    } else {
        http.Redirect(w, r, link, http.StatusNotFound)
        fmt.Println(err)
    }
}

func main () {

    s := &http.Server{
        Addr:           port,
    }

    // обработчики запросов:
    // - проверить, хочет юзер перейти по короткой ссылке или запрашивает фронтенд
    // сервиса
    http.HandleFunc("/", checkQuery)
    // - получить новую короткую ссылку
    http.HandleFunc("/reg_new_link/", registerNewLink)
    // - получить статистику переходов по короткой ссылке
    http.HandleFunc("/get_link_statistiсs/", getLinkStatistiсs)
    // открыть новую вкладку и перенаправить в нее вывод статистики посещений
    http.HandleFunc("/get_stat/", getStat)

    //запускаем сервер
    log.Fatal(s.ListenAndServe())
}
