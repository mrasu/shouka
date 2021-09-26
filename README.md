
# Shouka

Shouka is a tool to setup CI/CD environment easily.

## Why Shouka?

Compared with library or framework, companies don't change their way of delivery often because they're focusing on their business and don't have room for updating their development environment.  
Shouka helps you to change your continuous integration and continuous delivery by generating code for your environment.  
With generated code, you can set up a new way, evaluate it and import it to your code . 

https://user-images.githubusercontent.com/1549784/134805933-55e541d5-b2da-4b49-b8b9-fca8ffa1c90d.mp4

## How to use

Currently, Shouka supports AWS ECS that uses Docker image built by GitHub Actions.  
What you need to do are:

1. Download binary [here](https://github.com/mrasu/shouka/releases/tag/v0.0)
2. Run `shouka generate`
3. Answer questions by Shouka
4. Run `terraform init` and `terraform apply` under "terraforms" directory Shouka made
5. Push generated code to GitHub.
6. You will see GitHub Actions runs test, builds image and starts CodeDeploy to update a image in ECS cluster. 
