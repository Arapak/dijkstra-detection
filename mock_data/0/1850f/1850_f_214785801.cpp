// #pragma GCC optimize("Ofast")
// #pragma GCC target("sse,sse2,sse3,ssse3,sse4,popcnt,abm,mmx,avx,avx2,fma")
// #pragma GCC optimize("unroll-loops")
		
#include <bits/stdc++.h>
using namespace std;

typedef long long lli;
typedef long double ld;
typedef vector<int> vi;
typedef vector<ld> vld;
typedef vector<pair<lli , lli>> vpll;
typedef vector<pair<ld , ld>> vplld;
typedef pair<int,int> pii;
typedef vector<pair<int,int>> vpii;
typedef vector<lli> vl;
typedef pair<lli,lli> pll;
typedef priority_queue<lli> pq;
typedef priority_queue<pair<lli,lli>> pqp;

#define fi first
#define se second
#define pb push_back
#define mp make_pair

#define print(a) for(auto x:a) cout<<x<<" ";cout<<endl;
#define printarr(a , n) for(int i = 0 ; i < n  ;i ++) cout << a[i] << " "; cout << endl;
#define endl '\n'
#define sq(a) (a)*(a)
#define yes  cout << "YES" << endl;
#define no  cout << "NO" << endl;


#include <ext/pb_ds/assoc_container.hpp>
#include <ext/pb_ds/tree_policy.hpp>
using namespace __gnu_pbds;
  
#define ordered_set tree<int, null_type,less<int>,rb_tree_tag,tree_order_statistics_node_update>
#define inf 1e18


/*
	order_of_key(x) -> number of elements strictly smaller than x
	find_by_order(k) -> kth element 

	---------------------------------------------
	Think reverse -> graph problems
	Try converting into new operation
	check what happens to prefix sums
	graph trav , scc , topo , floyd ,  .... 
	copy code check
	index bounds
	maxsz
	set low , up bound be carefull
*/

/*

    ----------------------------------------------------------------------
	-------------- GO THROUGH CODE ONCE BEFORE SUBMITTING ----------------
    ----------------------------------------------------------------------

    ----------------------------------------------------------------------
	-------------- TLE : INFINITE LOOP , TIGHT CONSTRAINTS ---------------
	-------------- MLE : WIERD CONSTRAINT ARRAY SIZE       ---------------
	-------------- RE : POP EMPTY STACK QUEUE, INDEX BOUNDS --------------
    ----------------------------------------------------------------------

    Things get better by time, hope for the best ......
    
*/

void solve(){

	lli n;
	cin >> n;
	vl a(n);

	vi frq(n+1 , 0);
	vi act(n+1 , 0);
	
	for(int i = 0 ; i < n ; i ++)
	{
		cin >> a[i];
		if(a[i] <= n)
		{
			frq[a[i]]++;
		}
	}


	int ans = 0;

	for(int i = 1 ; i <= n ; i ++)
	{
		for(int j = i ; j <= n ; j += i)
		{
			act[j] += frq[i];
		}
	}


	for(int i  = 1 ; i <= n ; i ++) ans = max(ans , act[i]);

	cout << ans << endl;


}
int main(){

	ios::sync_with_stdio(0);cin.tie(0);cout.tie(0);
	

	#ifndef ONLINE_JUDGE
		freopen("input.txt", "r" , stdin);
		freopen("output.txt", "w" , stdout);
	#endif

	int t;cin>>t; while(t--)
    solve();
	
	return 0;
}

