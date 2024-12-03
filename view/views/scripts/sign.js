import PocketBase from 'https://cdn.jsdelivr.net/npm/pocketbase@0.22.1/+esm';

// 必要なリソース取得
const signInBtn = document.getElementById('sign-in-btn')  // ボタンを取得
const signOutBtn = document.getElementById('sign-out-btn')
const statusP = document.getElementById('status') // ステータス表示
const getAuthInfo = document.getElementById('get-auth-info')

// 値の設定
const pb = new PocketBase('http://127.0.0.1:8090');  // const pb = new PocketBase('https://pocketbase.io');
const provider = 'google';
///リスナー

// サインインボタン
signInBtn.addEventListener('click', async () => {
  try {

    // 認証
    const authData = await pb.collection('users').authWithOAuth2({ provider: 'google' });
    statusP.innerText = `Signin success! Signin as ${pb.authStore.record.email}`;

    // 値の表示
    console.log(pb.authStore.isValid);
    console.log(pb.authStore.token);
    console.log(pb.authStore.record.id);
    
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
