STEP 1: Pick your option (do this once for EACH app)
- You work on it by typing to an AI in Terminal on your Mac → Option A
- You work on it in your browser at replit.com → Option B
- You work on it on another website (netlify, bolt.new, lovable, etc.) → Option C
Not sure? Send me a screenshot of what you see when you work on it.

STEP 2: Paste the prompt for your option into the AI you used to build the app, and follow what it tells you. Being asked to log in or click "Authorize" is normal.

STEP 3: Text me the GitHub link it gives you.

Stuck at any point? Send me a screenshot, or we can FaceTime and you can share your screen.


=== OPTION A: Terminal on your Mac ===
Start your AI tool in Terminal the way you normally do for this app, then paste everything between the lines.

Heads-up: if your Mac pops up a box about installing "command line developer tools," click Install. If Terminal asks for your password, type your Mac password and press Enter. Nothing shows on screen while you type, and that's normal.

----- copy from here -----
I'm not a programmer, so please explain everything in plain language and walk me through anything I need to do, one step at a time.

I want to share this app's code with my friend (GitHub username: cwfreelance) in a private GitHub repository on my account. If we're not already in this app's project folder, help me find it first.

1. Install the GitHub CLI (gh) if needed, and help me log in to my GitHub account with it.
2. Set up git for this project if it isn't already.
3. Make sure no secrets get uploaded: add a .gitignore covering .env files, node_modules, and build/cache folders; create a .env.example with variable names only; and check the code AND the git history for hardcoded API keys or passwords. Don't read or print any secret values. If you find a secret, stop and explain how I should replace it before continuing.
4. Create a private GitHub repository and push the code.
5. Add cwfreelance as a collaborator on the repository.
6. Give me the repository link.
----- copy to here -----


=== OPTION B: Replit ===
Open the app in Replit and paste this into the Agent chat:

----- copy from here -----
I'm not a programmer, so please explain in plain language and guide me through any clicks one step at a time.

I want to share this app's code with my friend in a private GitHub repository on my account.
1. Help me connect this project to GitHub using the Git tool (Tools → Git → Connect to GitHub) and create a new private repository.
2. Push the current code, keeping the Replit config files. Don't include secrets, .env files with real values, database dumps, or dependency/cache folders.
3. Give me the repository link.
Don't change how the app is published.
----- copy to here -----


=== OPTION C: Netlify or another website ===
Paste this into the AI chat on the site where you built the app:

----- copy from here -----
I'm not a programmer, so please explain in plain language and guide me through any clicks one step at a time.

I want to share this app's original source code with my friend in a private GitHub repository on my account.
1. Help me connect this project to GitHub and create a new private repository. (On Netlify this is under Project configuration → Build & deploy → Continuous deployment → Repository → Link to repository → Push to new repository.)
2. Check that the repository has the editable source code, not just the built website files. If it only has built files, help me find the original source.
3. Make sure no secrets, .env files with real values, or database dumps get uploaded.
4. Tell me in plain language whether anything about how my live site gets published will change.
5. Give me the repository link.
----- copy to here -----


=== LAST STEP for Options B and C: invite me ===
On your new repository's page on github.com:
1. Click Settings (top of the page)
2. Click Collaborators in the left menu (GitHub may ask for your password)
3. Click "Add people," type cwfreelance, and click Add
Then text me the link!