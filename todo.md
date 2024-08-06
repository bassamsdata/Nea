# Short term (0.1.0) All Done 
- [x] Modify the updateVersionInfo function to have sort and and better unique number 
- [x] check the rollback version again
   - [x] get the version
   - [x] use the verison
   - [x] update the json file or we don't need to
- [x] add RemoveNightlyVersion
- [x] add feature to rollback based on date
- [x] add config file
- [x] Sort Function
    - [x] add sort function
    - [x] implement sort function

# before v0.1.0
- [x] Add update function
- [x] add config function to the setup function
- [ ] add readme
- [ ] add logging mechanism
- [x] add check for setup in install functions

## Current Bugs:
- [x] `clean nightly` this clean the oldest nightly
    - [x] so now it deletes the version but it doesn't update the config file
- [x] `clean stable `doesn't work
- [x] `clean stable all `deletes the stable dir as well, os ls local doesn't work
- [x] `use stable` doesn't work
- [x] use function for stable to version


## Testing
- [x] install nightly
- [x] install specific version
- [x] rollback
- [x] list
- [x] clean
- [x] use

## Long term
- [ ] use lipgloss for ui and colors - it didn't work
- [ ] add test cases
- [ ] able to pin a nightly version to never been deleted
- [ ] customize the location of the version manager
- [ ] add support for linux and windows
- [ ] add `nvm clean n` to clean the oldest n versions
    - like `nvm clean 5` to clean the oldest 5 versions

