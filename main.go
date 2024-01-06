package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/go-resty/resty/v2"
)

const (
	apiEndpoint = "https://api.openai.com/v1/chat/completions"
)

func getOpenAIResponse(apiKey, question string) (string, error) {
	client := resty.New()

	response, err := client.R().
		SetAuthToken(apiKey).
		SetHeader("Content-Type", "application/json").
		SetBody(map[string]interface{}{
			"model":      "gpt-3.5-turbo",
			"messages":   []interface{}{map[string]interface{}{"role": "system", "content": "Hi can you tell me " + question + "?"}},
			"max_tokens": 50,
		}).
		Post(apiEndpoint)

	if err != nil {
		return "", fmt.Errorf("error while sending the request: %v", err)
	}

	body := response.Body()

	var data map[string]interface{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return "", fmt.Errorf("error while decoding JSON response: %v", err)
	}

	content := data["choices"].([]interface{})[0].(map[string]interface{})["message"].(map[string]interface{})["content"].(string)
	return content, nil
}

type Post struct {
	Content string
	Date    string
}

type SocialGraph struct {
	Adj              map[string]map[string]struct{}
	UserContentTopic map[string]string
}

func NewSocialGraph() *SocialGraph {
	return &SocialGraph{
		Adj:              make(map[string]map[string]struct{}),
		UserContentTopic: make(map[string]string),
	}
}

func (g *SocialGraph) addUser(user string) {
	if _, exists := g.Adj[user]; !exists {
		g.Adj[user] = make(map[string]struct{})
		fmt.Printf("User '%s' added successfully!\n", user)
	} else {
		fmt.Printf("User '%s' already exists!\n", user)
	}
}

func (g *SocialGraph) addEdge(follower, following string) {
	if _, existsFollower := g.Adj[follower]; existsFollower {
		if _, existsFollowing := g.Adj[following]; existsFollowing {
			g.Adj[follower][following] = struct{}{}
			fmt.Printf("'%s' is now following '%s'.\n", follower, following)
		} else {
			fmt.Println("Error: Following user does not exist.")
		}
	} else {
		fmt.Println("Error: Follower user does not exist.")
	}
}

func (g *SocialGraph) printAdjList() {
	fmt.Println("Social Graph Adjacency List:")
	for user, followers := range g.Adj {
		fmt.Printf("%s follows: ", user)
		for follower := range followers {
			fmt.Printf("%s ", follower)
		}
		fmt.Println()
	}
}

func (g *SocialGraph) kosarajuDFS(user string, visited map[string]struct{}, s *stack) {
	visited[user] = struct{}{}

	for neighbor := range g.Adj[user] {
		if _, isVisited := visited[neighbor]; !isVisited {
			g.kosarajuDFS(neighbor, visited, s)
		}
	}

	s.Push(user)
}

func (g *SocialGraph) kosarajuDFSReverse(user string, visited map[string]struct{}, component *[]string) {
	visited[user] = struct{}{}
	*component = append(*component, user)

	for neighbor := range g.Adj[user] {
		if _, isVisited := visited[neighbor]; !isVisited {
			g.kosarajuDFSReverse(neighbor, visited, component)
		}
	}
}

func (g *SocialGraph) findSCC() {
	s := newStack()
	visited := make(map[string]struct{})

	for user := range g.Adj {
		if _, isVisited := visited[user]; !isVisited {
			g.kosarajuDFS(user, visited, s)
		}
	}

	adjTranspose := make(map[string]map[string]struct{})
	for user, followers := range g.Adj {
		for follower := range followers {
			if _, exists := adjTranspose[follower]; !exists {
				adjTranspose[follower] = make(map[string]struct{})
			}
			adjTranspose[follower][user] = struct{}{}
		}
	}

	visited = make(map[string]struct{})

	for !s.isEmpty() {
		user := s.Pop()

		if _, isVisited := visited[user]; !isVisited {
			var component []string
			g.kosarajuDFSReverse(user, visited, &component)

			fmt.Print("SCC: ")
			for _, u := range component {
				fmt.Printf("%s ", u)
			}
			fmt.Println()
		}
	}
}

func main() {
	socialGraph := NewSocialGraph()

	socialGraph.addUser("Alice")
	socialGraph.addUser("Bob")
	socialGraph.addUser("Charlie")
	socialGraph.addUser("David")
	socialGraph.addUser("Eva")

	socialGraph.addEdge("Alice", "Bob")
	socialGraph.addEdge("Bob", "Charlie")
	socialGraph.addEdge("Charlie", "Alice")
	socialGraph.addEdge("David", "Eva")
	socialGraph.addEdge("Eva", "David")

	socialGraph.UserContentTopic["Alice"] = "Cars"
	socialGraph.UserContentTopic["Bob"] = "BMW"
	socialGraph.UserContentTopic["Charlie"] = "Engine"
	socialGraph.UserContentTopic["David"] = "RPM"
	socialGraph.UserContentTopic["Eva"] = "Automobiles"

	socialGraph.printAdjList()
	socialGraph.findSCC()

	apiKey := "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"

	question := "what is the factorial of 10"
	response, err := getOpenAIResponse(apiKey, question)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(response)
}

type stack struct {
	items []string
}

func newStack() *stack {
	return &stack{items: make([]string, 0)}
}

func (s *stack) Push(item string) {
	s.items = append(s.items, item)
}

func (s *stack) Pop() string {
	if len(s.items) == 0 {
		return ""
	}
	item := s.items[len(s.items)-1]
	s.items = s.items[:len(s.items)-1]
	return item
}

func (s *stack) isEmpty() bool {
	return len(s.items) == 0
}
