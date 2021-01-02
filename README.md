dl-newcomic
===========

Tool to scape site *newcomic.info*.


## Install
``` shell
> go get -u github.com/cjun714/dl
```

## How to Use
``` shell
# download pages from 'newcomic.info/page/201' - newcomic.info/page/300',
# including index pages, detail pages and cover images.
# A direcotry /201-300/ will be created,
# index pages will be stored in ./201-300/,
# detail pages will be stored in ./201-300/pages/,
# cover images will be stored in ./201-300/images/.
> dl-newcomic 201 300
```

**Not**e: if lots of file downloading failed, perhaps it's caused by poor
performance of server of newcomic.info, just increasing `time.Sleep` value, 700
is always a save value.
