<center> 
# codeword 
</center>


## What Is It?

__codeword__ opens
a device, typically a disk drive and searches the contents for target strings.
It was originally written to look for unencrypted [classified data][1] that may have
unintentionally "leaked" onto unclassified hard drives.  

### Installation 

If go is installed:

```
go get github.com/hotei/codeword
```

### Features

* Simple - under 300 LOC - easy to customize
* Table driven - trivial to add/delete/change targets. 
 * [Wiki article][1] has some typical government (US - NATO etc) targets 
 * A commercial enterprise could 
search for their "secret project keywords" like "Project Unobtanium".  
* Fast - typically runs at near max HD read speed.  However, it will still take hours 
per terabyte.  1e12 / 1e8 = 10,000 seconds or about 3 hours at 100MB/sec.
* Shows progress bar during scan

### Limitations

* __codeword__ is trivially defeated by even the simplest encryption - such as xor.
The underlying assumption is that leaks are _unintentional_
* __codeword__ matches stop at disk block boundaries, so small OS disk block size can
mask matches that would otherwise be hits - even if the blocks are contiguous.
* __codeword__ is prone to false positives if given poorly chosen targets.  Some
examples of noise generating targets are mentioned in source.
* <font color="red">Because it opens the raw device - such as /dev/sda - 
you must be root to run this program.</font>
* targets are literals joined with "OR", you can't look for "TOP" AND "SECRET" AND "ULTRA"
 however this is easy to change in scanBlock()
Despite these weaknesses in method - in actual use it proved to be a valuable first check.  

### Usage

Typical usage is :

```codeword | tee codeword-results.txt```

### 
### TO DO
* Essential:
 * TBD
* Nice to have:
 * Disk to be checked should probably be a flag -disk="/dev/sda"
* Nice, but No Immediate Rqmt
 * Option for RE2 as matching mechanism 

### Change Log
* 2014-02-08 Enhance fold ASCII > 127 if requested (default is false)
* 2014-02-08 Store/Display hits only once per run.  map hits[sha256]int  if new hit same as
one in map just increment hitcount, dont redisplay

Comments can be sent to <hotei1352@gmail.com> or to user "hotei" at github.com.
License is BSD-two-clause, in file "LICENSE"

### Resources

* [Wikipedia: Classified Information][1]
* [How the US DoD marks documents] [2]
* [go reference spec] [3]
* [go package docs] [4]
* [codeword][5] program docs

[1]: http://en.wikipedia.org/wiki/Classified_information "http://en.wikipedia.org/wiki/Classified_information"
[2]: http://www.dtic.mil/dtic/pdf/customer/STINFOdata/DoD5200_1ph.pdf "http://www.dtic.mil/dtic/pdf/customer/STINFOdata/DoD5200_1ph.pdf"
[3]: http://golang.org/ref/spec/ "go reference spec"
[4]: http://golang.org/pkg/ "go package docs"
[5]: http://github.com/hotei/codeword "github.com/hotei/codeword"

 
> Redistribution and use in source and binary forms, with or without modification, are
> permitted provided that the following conditions are met:
> 
>    1. Redistributions of source code must retain the above copyright notice, this list of
>       conditions and the following disclaimer.
> 
>    2. Redistributions in binary form must reproduce the above copyright notice, this list
>       of conditions and the following disclaimer in the documentation and/or other materials
>       provided with the distribution.
> 
> THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDER ``AS IS'' AND ANY EXPRESS OR IMPLIED
> WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND
> FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL <COPYRIGHT HOLDER> OR
> CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR
> CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
> SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON
> ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING
> NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF
> ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

Documentation (c) 2015 David Rook 

