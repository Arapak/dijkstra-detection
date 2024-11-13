#include <bits/stdc++.h>

#define pb push_back
#define int long long
typedef long long ll;
#define ld long double
#define pai pair<int, int>
#define fr first
#define sc second
#define all(x) x.begin(), (x).end()
#define rall(x) x.rbegin(), (x).rend()
using namespace std;
const int INF = (int) 1e18;

void solve() {
  int n;
  cin >> n;
  vector<int> a(n);
  for (int i = 0; i < n; i++) {
    cin >> a[i];
  }
  vector<int> cnt(n + 2);
  for (int i = 0; i < n; i++) {
    a[i] = min(a[i], n + 1);
    cnt[a[i]]++;
  }
  vector<int> ans(n + 1);
  for (int i = 1; i <= n; i++) {
    for (int j = i; j <= n; j += i) {
      ans[j] += cnt[i];
    }
  }
  cout << *max_element(all(ans)) << '\n';
}

signed main() {
  ios::sync_with_stdio(false);
  cin.tie(nullptr);
  cout.tie(nullptr);
  int t;
  cin >> t;
  while (t--) {
    solve();
  }
}