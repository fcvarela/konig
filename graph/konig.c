#include <stdlib.h>
#include <stdio.h>

#include <GL/glew.h>
#include <GLFW/glfw3.h>

GLFWwindow *window;

static void error_callback(int error, const char* description) {
  fprintf(stderr, "Error: %s\n", description);
  exit(EXIT_FAILURE);
}

void init(int width, int height, int fullscreen) {
  fprintf(stderr, "Initializing...\n");

  if (!glfwInit()) {
    exit(EXIT_FAILURE);
  }

  glfwSetErrorCallback(error_callback);
  glfwWindowHint(GLFW_CONTEXT_VERSION_MAJOR, 3);
  glfwWindowHint(GLFW_CONTEXT_VERSION_MINOR, 3);
  glfwWindowHint(GLFW_OPENGL_FORWARD_COMPAT, 1);
  glfwWindowHint(GLFW_OPENGL_PROFILE, GLFW_OPENGL_CORE_PROFILE);

  if (fullscreen == 1) {
    // ignore requested width and height, we're going fullscreen!
    const GLFWvidmode *vidMode = glfwGetVideoMode(glfwGetPrimaryMonitor());
    window = glfwCreateWindow(vidMode->width, vidMode->height, "Konig", glfwGetPrimaryMonitor(), NULL);
  } else {
    window = glfwCreateWindow(width, height, "Konig", NULL, NULL);
  }
  if (!window) {
    glfwTerminate();
  }

  glfwMakeContextCurrent(window);
  glfwSwapInterval(1);
}

int update() {
  glfwSwapBuffers(window);
  glfwPollEvents();
  return glfwWindowShouldClose(window);
}

void shutdown() {
  fprintf(stderr, "Shutting down...\n");
  glfwDestroyWindow(window);
  glfwTerminate();
}
