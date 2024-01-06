# Markdown

Goldmark follows the [CommonMark](https://spec.commonmark.org/0.31.2/) markdown
specification.

## Goldmark

Goldmark ast fields/functions:

```
Heading
  Level   int  // depth
List
  Marker  byte // trigger char
  IsTight bool // false if list item are by empty lines
  Start   int  // (of ordered List)
ListItem
  Offset  int
Paragraph
TextBlock
Text
  Segment  segment // position in source text
  Text()   []byte  // return text
TaskCheckBox
  State
```

Example ast

```
Document
  Heading
    Text
  List
    ListItem
      Paragraph
        Text
      List
        ListItem
          TextBlock
            Text
        ListItem
          TextBlock
            Text
    ListItem
      Paragraph
        Text
      List
        ListItem
          TextBlock
            Text
```
