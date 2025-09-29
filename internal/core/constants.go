package core

const (
    // RepoDir is the name of the directory where all kitkat data is stored.
    RepoDir = ".kitkat"
    // ObjectsDir is the subdirectory for storing all content-addressable objects.
    ObjectsDir = ".kitkat/objects"
    // RefsDir is the subdirectory for storing references like heads and tags.
    RefsDir = ".kitkat/refs"
    // HeadsDir is the subdirectory for storing branch heads.
    HeadsDir = ".kitkat/refs/heads"
    // TagsDir is the subdirectory for storing tags.
    TagsDir = ".kitkat/refs/tags"
    // IndexPath is the full path to the index file.
    IndexPath = ".kitkat/index"
    // HeadPath is the full path to the HEAD file.
    HeadPath = ".kitkat/HEAD"
    // CommitsPath is the full path to the commit log file.
    CommitsPath = ".kitkat/commits.log"
)