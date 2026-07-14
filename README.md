# audioglyph

This tool is used to read barcodes of vinyl/CD records and use them to get more in depth metadata about the record using API calls to Discogs and MusicBrainz.

It will have a CLI tool to quickly insert a barcode to get metadata, and a WEB GUI tool to scan barcodes using handheld device and display the metadata in a user-friendly way.

There shall be a WEB GUI to display all records that has been scanned and few extra fields like if in collection or not. Would be fun to have this as a library of sort to keep track what is in the collection and what is not.
Main purpose is though to then integrate trackers such as RED or Orpheus to use their API to see if an album exists in their tracker so either skip getting it since we can already,
download it from a tracker. But if it is missing from the trackers it could be added and marked, perhaps even if it has a bounty.


Purpose first is to get the "initial" first, so database to store records and sort of barcode reader to get metadata from Discogs and MusicBrainz.

# Dependencies
- Cobra - https://github.com/spf13/cobra
- Gorm - https://github.com/go-gorm/gorm
- PostgreSQL - https://www.postgresql.org/