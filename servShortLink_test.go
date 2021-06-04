package main

import
(
    "net/http"
    "net/http/httptest"
    "testing"
    "io"
    "fmt"
)

// проверяет перенаправление с короткого урла на
func redirectTest( t *testing.T, link string, expectedStatus int) {
    req, err := http.NewRequest("GET", "/" + "?Url=" + link, nil)

    if err != nil {
        t.Fatal(err)
    }

    rr := httptest.NewRecorder()
    handler := http.HandlerFunc(redirect)

    handler.ServeHTTP(rr, req)

    if status := rr.Code; status != expectedStatus {
        t.Errorf("Возвращен неверный статус: получил %v ожидал %v",
            status, expectedStatus)
    }
}

func getCntGoogleCom( t *testing.T, expected string ) {
    req, err := http.NewRequest("GET", "/get_link_statistiсs/" + "?Url=" + baseUrl + "rGu2", nil)

    if err != nil {
        t.Fatal(err)
    }

    rr := httptest.NewRecorder()
    handler := http.HandlerFunc(getLinkStatistiсs)

    handler.ServeHTTP(rr, req)

    body, err := io.ReadAll(rr.Body)
    if err != nil {
        t.Fatal(err)
    }
    got := string(body)

    if expected != got {
        t.Errorf("ожидал: %v получил: %v",
            expected, got)
    }
}

// проверяет:
// - получение короткой ссылки для google.com
// - состояние счетчика посещений сразу после получения короткой ссылки (должно быть 0)
// - перенаправление с короткой ссылки на google.com
// - состояние счетчика посещений после перенаправления (должно быть 1 посещение)
func TestGoogleCom(t *testing.T) {
    req, err := http.NewRequest("GET", "/reg_new_link/" + "?Url=" + "https://www.google.com", nil)

    if err != nil {
        t.Fatal(err)
    }

    rr := httptest.NewRecorder()
    handler := http.HandlerFunc(registerNewLink)

    handler.ServeHTTP(rr, req)

    if status := rr.Code; status != http.StatusOK {
        t.Errorf("Возвращен неверный статус: получил %v ожидал %v",
            status, http.StatusOK)

    } else {
        body, err := io.ReadAll(rr.Body)

        if err != nil {
            t.Fatal(err)
        }
        expected := fmt.Sprintf("Короткая ссылка для https://www.google.com - %s \n", baseUrl + "rGu2")
        got := string(body)

        if expected != got {
            t.Errorf("ожидал: %v получил: %v",
                expected, got)
        }

        expected  = fmt.Sprintf("Адрес %s посещали 0 раз \n", baseUrl + "rGu2")
        getCntGoogleCom(t, expected)
        expected  = fmt.Sprintf("%s", baseUrl + "rGu2")
        redirectTest(t, expected, http.StatusSeeOther )
        expected  = fmt.Sprintf("Адрес %s посещали 1 раз \n", baseUrl + "rGu2")
        getCntGoogleCom(t, expected)
    }
}

// проверяет возвращаемую ошибку, есть дать пустой длинный урл
func TestEmptyUrl(t *testing.T) {
    req, err := http.NewRequest("GET", "/reg_new_link/", nil)
    if err != nil {
        t.Fatal(err)
    }

    rr := httptest.NewRecorder()
    handler := http.HandlerFunc(registerNewLink)

    handler.ServeHTTP(rr, req)

    if status := rr.Code; status != http.StatusOK {
        t.Errorf("Возвращен неверный статус: получил %v ожидал %v",
            status, http.StatusOK)

    } else {
        body, err := io.ReadAll(rr.Body)

        if err != nil {
            t.Fatal(err)
        }

        expected:= "Введенный url неполный! \n"
        got := string(body)

        if expected != got {
            t.Errorf("ожидал: %v получил: %v",
                expected, got)
        }
    }
}

// проверяет перенаправление на длинный урл при заданной несуществующей в системе
// короткой ссылке
func TestRedirectNotExistingLink( t *testing.T ) {
    redirectTest(t, "http://www.some_website.ru", http.StatusNotFound )
}
