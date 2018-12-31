# iscardcopied
check if a memory card is copied to disk

This tool can be used to verify if all photos of a memory card are copied to a folder
on a hard drive.

example usage:

    iscardcopied --disk=E: --card=G:

will check if all videos and images on memory card in drive G: are somewhere on E:

List of file extensions to care for is in main.go.

It will create a `uncopied.html` file containing links to all files not copied,
if option `--copy` is used, files can be copied at same time.

example usage:

    iscardcopied --disk=E: --card=G:  --copy=E:\MISC

After running this command it should be safe to format the card.

Do not blame me if anything goes wrong...
