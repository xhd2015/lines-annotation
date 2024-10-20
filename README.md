# Lines Annotation
This library abstracts a general representation of code, which braces line coverage, block annotation and diff info.

Then, built on top of this representation, coverage info from different commits are able be merged together, which facilitates coverage inheritance.

# Library Structure
The following packages define different operations on the annotation structure defined in package [model](./model/):
- [load](./load): load primitive annotations like diff changes
- [filter](./filter): filter based on criteria
- [map](./map): map from one form to another, for example, prepend path to files
- [compute](./compute): compute derived annotations from primitvie annotations and other derived annotations

# Why Lines Annoation?
I find there are dozens of representations of code coverage, e.g. xml in java, plain text in golang. They correspond to different language-specific structures, which makes it quite difficult to develop a general coverage review UI.

To address this issue, I've created this library, firstly used in golang, then expanded to js. It has been proven universal for common usage.

# Usage
This project is used by [xgo](https://github.com/xhd2015/xgo) as core library to provide coverage report.