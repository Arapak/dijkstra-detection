#include <iostream>
#include <vector>
#include <queue>
#include <climits>

using namespace std;

const long long INF = LLONG_MAX;

class Node
{
	public:
		int data = 0;
		long long dist = INF;
		vector<pair<int, long long> > edges;
		bool visited = false;
		int pre = 0;
};

class mycomp
{
	public:
		bool operator()(Node a, Node b)
		{
			return a.dist > b.dist;
		}
};

int main()
{
	int n, m;
	int a, b, c;
	cin >> n >> m;
	Node graph[100001];
	priority_queue<Node, vector<Node>, mycomp> pq;
	for (int i = 0; i < m; i++)
	{
		cin >> a >> b >> c;
		graph[a].edges.push_back(make_pair(b, c));
		graph[b].edges.push_back(make_pair(a, c));
		graph[a].data = a;
		graph[b].data = b;	
	}
	graph[1].dist = 0;
	pq.push(Node());
	pq.push(graph[1]);
	graph[1].visited = true;		
	Node curr = pq.top();
	pq.pop();
	long long d, e;
	while (!pq.empty())
	{
		a = curr.edges.size();
		for (int i = 0; i < a; i++)
		{
			d = graph[curr.edges[i].first].dist;
			e = curr.dist + curr.edges[i].second;
			if (d > e)
			{
				graph[curr.edges[i].first].dist = e;
				graph[curr.edges[i].first].pre = curr.data;
			}
			if (!graph[curr.edges[i].first].visited)
			pq.push(graph[curr.edges[i].first]);
		}
		for (int i = 0; i < n && pq.size() > 0; i++)
		{
			if (!graph[pq.top().data].visited)
			{
				curr = pq.top();
				break;
			}
			pq.pop();
		}
		graph[curr.data].visited = true;
		pq.pop();
	}
	if (graph[n].dist == INF)
	{
		cout << -1;
		return 0;
	}
	vector<int> path;
	a = n;
	while (a != 1)
	{
		path.push_back(a);
		a = graph[a].pre;
	}
	path.push_back(1);
	for (int i = path.size() - 1; i > 0; i--)
	{
		cout << path[i] << " ";		
	}
	cout << n;
	return 0;
}
