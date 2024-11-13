#include <bits/stdc++.h>
#define int long long

const int MAXN = 1000000 + 10;
const int MAXM = 5000 + 10;
const int INF = 1e18;

int t[MAXM][MAXM];
bool bl[MAXN];
char s[MAXM][MAXM];
int a[MAXN];
int b[MAXN];
int n, m, q, k, c; 
int x, y, z;
int tests;

void reset()
{
    
}

int vals[MAXN];
std::map <int,int> cnt;
void solve()
{
    cnt.clear();
    for (int i = 1 ; i <= n ; ++i)
    {
        cnt[a[i]]++;
        vals[i] = 0;
    }

    for (const auto &[key, value] : cnt)
    {
        int pos = 0;
        while (pos + key <= n)
        {
            pos += key;
            vals[pos] += value;
        }   
    }

    int max = 0;
    for (int i = 1 ; i <= n ; ++i)
    {
        max = std::max(max, vals[i]);
    }

    std::cout << max << '\n';
}

void input() 
{ 
    std::cin >> n;
    for (int i = 1 ; i <= n ; ++i)
    {
        std::cin >> a[i];
    }
}

void fastIOI() 
{
    std::ios_base :: sync_with_stdio(0);
    std::cout.tie(nullptr); 
    std::cin.tie(nullptr);
}

signed main () 
{
    fastIOI();    
    std::cin >> tests;

    while (tests--) 
    {    
        reset();
        input();
        solve();
    }
 
    return 0;
}