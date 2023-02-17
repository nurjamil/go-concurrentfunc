# Concurrent function runner

This package is to allow running 2 or more function to running concurrently using go routine and channel
## Example
```go
type People struct {
    Name string
}

type Animal struct {
    Type AnimalType
}

type AnimalType int

const (
	Type AnimalType = iota
	DOG
)

func (a AnimalType) String() string {
	switch a {
	case DOG:
		return "DOG"
	}

	return "UNKNOWN"
}

func main() {
    firstFunc := func(ctx context.Context) (any, error) {
        // mocking context checking on database or other dependency
        select {
        case <-ctx.Done():
            return nil, ctx.Err()
        default:
        }

        firstFuncRes := People{Name: "example"}

        return firstFuncRes, nil
    }

    secondFunc := func(ctx context.Context) (any, error) {
        select {
        case <-ctx.Done():
            return nil, ctx.Err()
        default:
        }

        secondFuncRes := Animal{Type: DOG}

        return secondFuncRes, nil
    }

    res, errs, err := concurrentfunc.Exec(context.Background(), time.Second, firstFunc, secondFunc)
    if err != nil {
        // root cause or the first appear error
        log.Error(err)
    }
    // convert the interface to your function type
    fmt.Println(res[0].(People))
    fmt.Println(res[1].(Animal))

    // checking error for each function
    fmt.Println(errs[0])
    fmt.Println(errs[1])


}


```
