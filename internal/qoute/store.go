package qoute

import (
	"math/rand"
)

var quotes = []string{
	"The only true wisdom is in knowing you know nothing.",
	"There are many ways of going forward, but only one way of standing still",
	"The journey of a thousand miles begins with one step",
	"If you talk to a man in a language he understands, that goes to his head. If you talk to him in his language, that goes to his heart",
	"It’s not what happens to you, but how you react to it that matters",
	"If you don’t know where you are going, any road will get you there",
	"The best thing to hold onto in life is each other",
	"Beware of false knowledge; it is more dangerous than ignorance",
	"You are a product of your environment. So choose the environment that will best develop you toward your objective. Analyze your life in terms of its environment. Are the things around you helping you toward success – or are they holding you back?",
	"Obstacles are those frightful things you see when you take your eyes off your goal",
	"For beautiful eyes, look for the good in others; for beautiful lips, speak only words of kindness; and for poise, walk with the knowledge that you are never alone",
	"The pessimist complains about the wind; the optimist expects it to change; the realist adjusts the sails",
	"In every walk with nature one receives far more than he seeks",
	"Discipline is the bridge between goals and accomplishment",
	"We are what our thoughts have made us; so take care about what you think. Words are secondary. Thoughts live; they travel far",
	"In three words I can sum up everything I’ve learned about life: it goes on",
	"The most common way people give up their power is by thinking they don’t have any",
	"Life is not a problem to be solved, but a reality to be experienced",
	"Opportunity is always knocking. The problem is that most people have the self-doubt station in their heads turned up way too loud to hear it",
	"The greatest day in your life and mine is when we take total responsibility for our attitudes. That’s the day we truly grow up",
}

type (
	// Store is a store of quotes
	Store struct {
		quotes []string
	}
	// Option is a store option
	Option func(*Store)
)

// WithQuotes sets quotes for a store
func WithQuotes(quotes []string) Option {
	return func(s *Store) {
		s.quotes = quotes
	}
}

// New returns a new quote store with predefined quotes
func New(opts ...Option) *Store {
	store := &Store{quotes: quotes}
	for _, opt := range opts {
		opt(store)
	}
	return store
}

// Random returns a random quote
func (q *Store) Random() string {
	return q.quotes[rand.Intn(len(q.quotes))]
}
