package hasher


func fnv1a(key string) uint64 {
    const (
        offset64 = 14695981039346656037
        prime64  = 1099511628211
    )
    h := uint64(offset64)
    for i := 0; i < len(key); i++ {
        h ^= uint64(key[i])
        h *= prime64
    }
    return h
}

func ShardFor(key string, n uint64) uint64 {
    return fnv1a(key) % n
}