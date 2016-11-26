#include <stdlib.h>
#include <stdio.h>

#include <GL/glew.h>

#if defined ( __APPLE__ )
#include <OpenCL/opencl.h>
#define GLFW_EXPOSE_NATIVE_COCOA
#define GLFW_EXPOSE_NATIVE_NSGL
#include <OpenGL/OpenGL.h>
#elif defined ( WIN32 )
#include <CL/cl.h>
#define GLFW_EXPOSE_NATIVE_WIN32
#define GLFW_EXPOSE_NATIVE_WGL
#else
#include <CL/cl.h>
#define GLFW_EXPOSE_NATIVE_X11
#define GLFW_EXPOSE_NATIVE_GLX
#endif

#include <GLFW/glfw3.h>
#include <GLFW/glfw3native.h>

GLFWwindow *window;

static void error_callback(int error, const char* description) {
  fprintf(stderr, "Error (%d): %s\n", error, description);
  exit(EXIT_FAILURE);
}

void graphview_init(int width, int height, int fullscreen) {
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

  glfwMakeContextCurrent(window);
  glfwSwapInterval(1);
}

int graphview_update() {
  glfwSwapBuffers(window);
  glfwPollEvents();
  return glfwWindowShouldClose(window);
}

void graphview_shutdown() {
  glfwDestroyWindow(window);
  glfwTerminate();
}
