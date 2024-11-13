#include <bits/stdc++.h>
#define INF 2e18
#define N 111111
using namespace std;
typedef long long ll;
typedef pair<ll,int> pli;
int n,m,p[N];
vector<pli> g[N];
ll d[N];
priority_queue<pli> q;
void trace(int k){if(p[k]!=-1) trace(p[k]); printf("%d ",k);}
int main(){
	scanf("%d%d",&n,&m);
	while(m--){
		int u,v,l;scanf("%d%d%d",&u,&v,&l);
		g[u].push_back(pli(l,v)),g[v].push_back(pli(l,u));
	}
	fill(p,p+n+5,-1); 
	fill(d,d+n+5,INF);
	q.push(pli(d[1]=0,1));
	while(!q.empty()){
		int u=q.top().second; q.pop();
		for(pli e: g[u]){
			int v=e.second; ll t=d[u]+e.first;
			if(t<d[v])p[v]=u,q.push(pli(-(d[v]=t),v));
		}
	}
	if(d[n]==INF) {printf("-1\n"); return 0;}
	trace(p[n]);
	printf("%d\n",n);
	return 0;
}