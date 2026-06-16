// Package upsilonserializer defines the canonical schema version for Upsilon engine state blobs.
//
// When the engine's persisted state shape changes (new fields, renamed keys, restructured
// sub-objects), this constant MUST be incremented. Resurrection guards in upsilonapi compare
// the embedded version against this value and refuse to hydrate incompatible blobs.
//
// Data-migration note: blobs persisted before this versioning scheme was introduced carry
// no serializer_version field (it will unmarshal as zero). Those blobs must be treated as
// stale and resurrection must be refused — the match should be concluded server-side and the
// game_state_cache column cleared. No automated DB migration is implemented here; that is a
// separate ops concern.
package upsilonserializer

// CurrentSerializerVersion is the schema version embedded into every BoardState blob.
// Increment this constant whenever the serialized engine state shape changes in a way that
// would make an older blob unsafe or incorrect to deserialize.
//
// History:
//
//	1 — initial versioned schema (WP-D2 / audit risk R7)
const CurrentSerializerVersion = 1
