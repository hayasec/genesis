# Genesis

Missed golang generic stdlib for slices and channels.

**Some functions:**

+ Filter, Map, Reduce.
+ Min, Max, Sum.
+ Permutations, Product.
+ Any, All.
+ Contains, Find.
+ Shuffle, Sort.
+ Range, Count, Cycle.

And much more.

**Features:**

+ Typesafe.
+ Sync and async versions.
+ For slices and channels.
+ Pre-generated for all built-in types.

```bash
go get github.com/life4/genesis
```

## Examples

Find minimal value in a slice of ints:

```go
s := []int{42, 7, 13}
min := genesis.SliceInt{s}.Min()
```

Double values in a slice of ints:

```go
s := []int{4, 8, 15, 16, 23, 42}
double := func(el int) int { return el * 2 }
doubled := genesis.SliceInt{s}.MapInt(double)
```

See [docs](./docs) to dive deeper.

## Custom types

Genesis contains pre-generated code for common built-in types. So, in most cases you can just use it. However, if you want to use genesis for custom types, things become a little bit more complicated. The first option is to use an empty interface. For example:

```go
type UserId int
ids := []UserId{1, 2, 3, 4, 5}
// https://github.com/golang/go/wiki/InterfaceSlice
idsInterface := make([]interface{}, len(ids), len(ids))
for i := range ids {
	idsInterface[i] = ids[i]
}
index := genesis.SliceInterface{idsInterface}.FindIndex(
	func(el interface{}) bool { return el.(UserId) == 3 },
)
fmt.Println(index)
// Output: 2
```

Another option is to generate genesis code for your own type.

## Generation

Install requirements

```bash
python3 -m pip install --user -r requirements.txt
```

Re-generate everything for built-in types:

```bash
python3 -m generate
```

Generate a new package with given types:

```bash
python3 -m generate
```

# [-> DOCUMENTATION <-](./docs)

![](./gopher.png)
