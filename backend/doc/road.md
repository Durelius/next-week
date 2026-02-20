### init 26/02/17

Studying generic for structs in Golang (https://go.dev/doc/tutorial/generics, https://stackoverflow.com/questions/68166558/generic-structs-with-go)



### 26/02/19 implementing AVL according to w3schools python implementation
https://www.w3schools.com/python/python_dsa_avltrees.asp

Generating test suite from this prompt in Claude sonnet 4.6, Gemini 3 'Thinking' and ChatGPT free:

Generate an Go test suite for a custom AVL node, that will be placed in an external 'test' package.
You can use these methods from the avl package i created:
New(data of type cmp.ordered), returns new node
Insert(data of type cmp.ordered) returns root node
String() returns comma and space separated list of all nodes in order
Data() returns the data of the current node
The AVL tree currently does not support duplicates, so do not include duplicates in the test suite
Make sure the testing covers all edge cases and has a lot of data.
You don't have access to the actual node class, it is private. You can only use the methods above.


Gemini still tried to use the node struct, thinking it was public, so I manually adjusted it.


The test suites are named after the AI who generated them.


Shortly after this I added a key, value structure to the AVL tree. So only the key needs to be comparable, and the value can be whatever. Need to redo the tests though with this in mind.
Also made the values be slices, so we get a slice of values for the duplicates. 

### 26/02/20 Time for additional testing

Now we have changed the whole AVL to use keys and slice values, and added a tree struct with public methods. It's time to try and make our dear LLMs generate test suites again.
Generating test suite from this prompt in Claude sonnet 4.6, Gemini 3 'Thinking' and ChatGPT free:
Prompt:
//promptstart
This is a custom AVL tree, built by me. The tree struct has these public methods:


func New\[K cmp.Ordered, V any]() *Tree\[K, V] {
func (t *tree \[K, V]) Insert(key K, value V) {
func (t *tree\[K, V]) Delete(key K) {
func (t *tree\[K, V]) Height() int {
func (t *tree\[K, V]) Size() int {
func (t *tree\[K, V]) Min() (K, []V, bool) {
func (t *tree\[K, V]) Find(key K) ([]V, bool) {
func (t *tree \[K, V]) Contains(key K) bool {
func (t *tree\[K, V]) String() string {
func (t *tree\[K, V]) Print() {

Use only these public methods, from an external package, to create an extensive test suite that covers all usual edge cases and weird behaivours.
Use a LOT of data.

//promptend


