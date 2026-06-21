package obs_test

import (
	"context"
	"os/exec"
	"testing"

	"github.com/omcrgnt/demo/v2/pkg/obs"
	"github.com/omcrgnt/demo/v2/pkg/obs/testdata"
)

func TestObserveInAction(t *testing.T) {
	cmd := exec.Command("go", "generate", "./testdata/repo.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("generate failed: %v\noutput: %s", err, out)
	}

	rawRepo := testdata.NewRepo("card-storage")

	observer, ok := any(rawRepo).(obs.Observer)
	if !ok {
		t.Fatal("resource does not implement obs.Observer")
	}

	// 4. Получаем "наблюдаемую" версию
	// Она возвращает any, поэтому приводим к интерфейсу, который нужен сервису
	decorated := observer.Observe().(interface {
		Save(context.Context, string) error
	})

	// 5. Проверяем работу декоратора (оборачивание ошибок)
	err := decorated.Save(context.Background(), "") // вызываем ошибку
	if err == nil {
		t.Fatal("expected error, but got nil")
	}

	// Проверяем формат: [Label] Method: Error
	expected := "[card-storage] Save: empty data"
	if err.Error() != expected {
		t.Errorf("wrong error format\nwant: %s\ngot:  %s", expected, err.Error())
	}
}
