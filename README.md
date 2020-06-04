
# css3fmt

[css3fmt](https://blitiri.com.ar/git/r/css3fmt) is an auto-formatter for
[CSS](https://en.wikipedia.org/wiki/Cascading_Style_Sheets) files.

It is not particularly fancy or smart, but it is simple and can automatically
format most CSS files.


## Install

css3fmt is written in Go.

```sh
go get blitiri.com.ar/go/css3fmt
```


## Editor integration

### vim

Put the following into your `.vimrc` file to auto-indent on save:

```vim
function! CSSFormatBuffer()
        let l:curw = winsaveview()
        let l:tmpname = tempname()
        call writefile(getline(1,'$'), l:tmpname)
        let l:out = system("css3fmt " . l:tmpname) 
        call delete(l:tmpname)  
        if v:shell_error == 0           
                try | silent undojoin | catch | endtry
                silent %!css3fmt     
        else    
                echoerr l:out
        endif
        call winrestview(l:curw)
        return v:shell_error == 0
endfunction
autocmd filetype css
  \ autocmd bufwritepre <buffer> call CSSFormatBuffer()
```


## Contact

If you have any questions, comments or patches please send them to
albertito@blitiri.com.ar.

