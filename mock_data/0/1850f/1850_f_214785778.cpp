#include<bits/stdc++.h>
using namespace std;

#ifdef LOCAL
    #include"debug.h"
#else
    #define dbg(...) 
#endif

#define all(x) begin(x),end(x)
#define int long long

int32_t main(){
	ios_base::sync_with_stdio(0);
	cin.tie(0);
	
	int tt;
	cin >> tt;
	while(tt--){
		int n;
		cin >> n;
		vector<int > a(n);
		map<int,int> f;
		for(int i =0 ;i<n;i++){
			cin >> a[i];
			f[a[i]]++;
		}
		vector<int > val(n + 1);
		for(auto x : f){
			// they jump with x.first and are x.second
			
			for(int i = x.first ;i<=n;i+=x.first){
				val[i] += x.second;
			}
		}
		cout << *max_element(all(val)) << endl;
	}

	return 0;
}