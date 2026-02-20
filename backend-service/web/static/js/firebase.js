import { initializeApp } from "https://www.gstatic.com/firebasejs/10.7.1/firebase-app.js";
import { getAuth } from "https://www.gstatic.com/firebasejs/10.7.1/firebase-auth.js";

const firebaseConfig = {
  apiKey: "AIzaSyC-sTo64NctAETMidi26kbB2mRpBKmmFLI",
  authDomain: "readers-app-87268.firebaseapp.com",
  projectId: "readers-app-87268",
};

export const app = initializeApp(firebaseConfig);
export const auth = getAuth(app);
