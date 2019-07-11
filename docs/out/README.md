# This repo hosts sample MD template for API creation and Data Dictionary YAML creation. 

## API Template Notes
* Refer to template-post.md and template-get.md files 

Below is a reference for the Markdown syntax used in Slate.

Headers
For headers:

# Level 1 Header
## Level 2 Header
### Level 3 Header
Note that only level 1 and 2 headers will appear in the table of contents.

Paragraph Text
For normal text, just type your paragraph on a single line.

This is some paragraph text. Exciting, no?
Make sure the lines above and below your paragraph are empty.

Code Samples
For code samples:
	```json
	# Some JSON code!
	```
Code samples will appear in the dark area to the right of the main text. Slate recommends positioning code samples right under headers in your markdown file.

Code Annotations
For code annotations:

> This is a code annotation. It will appear in the area to the right, next to the code samples.
Code annotations are essentially the same thing as paragraphs, but they'll appear in the area to the right along with your code samples.

Formatting code annotations (right column) with its explanation (center column)
In order to correctly format the code in the right column with the text in the center column the code snippet should go first, e.g.

    # My Title

    ## My Subtitle

    ```java
        Code snippet
    ```

    ```json
    #Code snippet
    ```

    Whatever that goes in the center column.
Putting the content for the center column first and the code snippet afterwards will cause the code snippet to be aligned with the last line of the center column, if the code snippet goes first, then the code is aligned with the first line of the center column.

Tables
Slate uses PHP Markdown Extra style tables:

Table Header 1 | Table Header 2 | Table Header 3
-------------- | -------------- | --------------
Row 1 col 1 | Row 1 col 2 | Row 1 col 3
Row 2 col 1 | Row 2 col 2 | Row 2 col 3
Note that the pipes do not need to line up with each other on each line.

If you don't like that syntax, feel free to use normal html <table>s directly in your markdown.

Formatted Text
This text is **bold**, this is *italic*, this is an `inline code block`.
You can use those formatting rules in tables, paragraphs, lists, wherever, although they'll appear verbatim in code blocks.

Lists
1. This
2. Is
3. An
4. Ordered
5. List

* This
* Is
* A
* Bullet
* List
Links
This is an [internal link](#error-code-definitions), this is an [external link](http://deltadental.com).
Notes and Warnings
You can add little highlighted warnings and notes with just a little HTML embedded in your markdown document:

<aside class="notice">
You must replace `meowmeowmeow` with your personal API key.
</aside>
Use class="notice" for blue notes, class="warning" for red warnings, and class="success" for green notes.
