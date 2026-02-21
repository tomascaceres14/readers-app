import {
  createUserWithEmailAndPassword,
  signInWithPopup,
  GoogleAuthProvider,
} from "https://www.gstatic.com/firebasejs/10.7.1/firebase-auth.js";
import { auth } from "./firebase.js";

document
  .getElementById("register-form")
  .addEventListener("submit", async (event) => {
    event.preventDefault();
    const form = event.target;
    const email = form.email.value;
    const password = form.password.value;
    try {
      const userCredential = await createUserWithEmailAndPassword(
        auth,
        email,
        password,
      );
      const idToken = await userCredential.user.getIdToken();
      form.password.value = "";
      form.token.value = idToken;
      form.submit();
    } catch (error) {
      console.error(error);
    }
  });

document
  .getElementById("register-google")
  .addEventListener("submit", async (event) => {
    event.preventDefault();
    const form = event.target;
    try {
      const provider = new GoogleAuthProvider();
      const cred = await signInWithPopup(auth, provider);
      const idToken = await cred.user.getIdToken();
      form.token.value = idToken;
      form.submit();
    } catch (error) {
      console.log(error);
    }
  });
