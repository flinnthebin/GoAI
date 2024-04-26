# iago                                                                                                                                                                                                          
iago - a command-line utility for accessing openAI LLMs

# Note : OpenAI have aggressively refactored the means by which their LLMs define "valid JSON"
This version is officially deprecated until we write a custom JSON tool to accompany iago

## File Interaction                                                                                                                                                                                                
iago is an extension of shGPT, designed to improve reliability when handling large files.

iago also seems to run slightly faster than the original shGPT bash script.

## creating 

```
git clone https://github.com/flinnthebin/iago

make run

mv iago /usr/bin/

```
## init.vim

To use :iago to launch the script from your current vim instance.

```
command iago call iago()

function! iago()
    let shgpt_path = expand('~/usr/bin/iago')
    let tokens = 4096
    let temperature = 0.1
    let current_file = expand('%:p')
    let creator_ftype = fnamemodify(current_file, ':e')
    let outfile = 'output.' . creator_ftype
    let shell_cmd = '! ' . shellescape(shgpt_path) . ' -s ' . tokens . ' -t ' . temperature . ' ' . shellescape(current_file)
    execute shell_cmd
    execute 'vsplit ' . outfile
endfunction
```

## Tokens

A token is considered to be 4 characters. The number of tokens set must be enough to include the length
of the provided input and expected output.

I'm pretty sure the hard limit for GPT-4 is 4096 currently?

## Temperature

The temperature parameter controls the randomness of the output. 

My personal perference is 0.1, feel free to adjust this.

Low (0.1 - 0.3)
Predictable, deterministic output.

Medium (0.4 - 0.6)
An adaptable range balancing creativity and predictability

High (0.7 - 1.0)
Creative and unpredictable. Strong possibility of hallucination.



