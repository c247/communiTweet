# communiTweet

### 1
- Cool idea I had to find communities in Social Media. Assuming that you have a social graph where each node is a user and each directed edge between user a and user b represents a follower relationship, we can determine Strongly Connected Components (SCCs) utilizing Kosaraju's algorithm. This can be used to suggest users to be a part of a community with others who have similar interests. Assuming that two users being in the same strongly connected component implies that they enjoy the same type of content.

### 2
- Next, when looking at posts from a given community, I determine the dependencies in the posts using AI (example: posts related to engines should come after posts related to automobiles - to help make the scrolling experience for the user more story-like). This is done by creating a Directed Acyclic Graph based on the dependencies in the posts and displaying the posts in a topological order.
