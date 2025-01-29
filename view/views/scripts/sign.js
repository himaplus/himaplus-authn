import PocketBase from 'https://cdn.jsdelivr.net/npm/pocketbase@0.22.1/+esm'; // https://www.jsdelivr.com/package/npm/pocketbase

// 必要なリソース取得
const signInBtn = document.getElementById('sign-in-btn')  // ボタンを取得
const signOutBtn = document.getElementById('sign-out-btn')
const statusP = document.getElementById('status') // ステータス表示
const getAuthInfo = document.getElementById('get-auth-info')

// 値の設定
const pbUrl = `${location.protocol}//${location.hostname}${location.port != ""? `:${location.port}`: ""}`; // const pbUrl = `${location.protocol}//${location.hostname}` + (location.port != "" ? `:${location.port}` : "");
const pb = new PocketBase(pbUrl);  // const pb = new PocketBase('https://pocketbase.io');
const provider = 'google';

//リスナー

// サインインボタン
signInBtn.addEventListener('click', async () => {
  try {
    // 認証
    const authData = await pb.collection('users').authWithOAuth2({
      provider: provider,
      scopes: [ // The presence or absence of this place may be irrelevant.
        "https://www.googleapis.com/auth/calendar",
        "https://www.googleapis.com/auth/calendar.readonly",
        "https://www.googleapis.com/auth/userinfo.email",
        "https://www.googleapis.com/auth/userinfo.profile",
      ],
      urlCallback: async (url) => {
        const adjustedUrl = new URL(url);
        console.log(adjustedUrl);
  
        // クエリを追加
        adjustedUrl.searchParams.set('access_type', 'offline')
        adjustedUrl.searchParams.set('prompt', 'consent');

        console.log(adjustedUrl.toString());

        // 新しいタブで同意画面
        window.open(adjustedUrl, '_blank');  // 第二引数を省略すると同じタブで開く      
      }
    });
    statusP.innerText = `Signin success! Signin as ${pb.authStore.record.email}`;

    // 値の表示
    console.log("pb.authStroe.isValid: ", pb.authStore.isValid);
    console.log("pb.authStore.token: ", pb.authStore.token);
    console.log("pb.authStore.record.id: ", pb.authStore.record.id);
  } catch (err) {
    console.error(err);
    statusP.innerText = 'Signin failed!';
  }
})

// サインアウトボタン
signOutBtn.addEventListener('click', () => {
  try {
    // 認証
    pb.authStore.clear();
    statusP.innerText = "Signout success!";
  } catch (err) {
    console.error(err);
    statusP.innerText = 'Signout failed!';
  }
})

// 認証情報取得
getAuthInfo.addEventListener('click', async () => {
  try {
    // データの確認
    console.log(pb.authStore.isValid);
    console.log(pb.authStore.token);
    console.log(pb.authStore.model.id);

    console.log(pb.authStore);
  } catch (err) {
    console.error(err);
  }
})
