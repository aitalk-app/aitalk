# AI Talk CLI

[中文说明](/README-cn.md)

Create and share AI talks https://ai-talk.app

## Usage
### To create a talk with two AI automatically

```bash
aitalk --topic "Is PHP the best programming language?" --role "A C++ programmer who think C++ is best" --role "A PHP programmer who believe PHP is best"
```

#### Use a different language

```bash
# show supported languages
aitalk lang

# specify language
aitalk --lang {lang} --topic "Is PHP the best programming language?" --role "A C++ programmer who think C++ is best" --role "A PHP programmer who believe PHP is best"
```

### To create a talk with AI interactively

```bash
aitalk --topic "Will AI replace human?"
```

## Manage your talks (optional)
If you want to manage your talks on https://ai-talk.app, you need to authenticate. Run the following command and open the displayed URL in your web browser:

```bash
aitalk auth
```

If you skip this step now, you can run the above command later and all previously created talks will automatically get assigned to your profile.
