#include <iostream>
#include <fstream>
#include <set>
#include <filesystem>

std::set<std::string> get_songs_from_file(std::string path) {
  std::ifstream fs { path, std::fstream::in };
  std::set<std::string> res {};
  if (fs.is_open()) {
    while (fs.good()) {
      std::string line;
      getline(fs, line);
      res.insert(line);
    }
    fs.close();
  } else {
    throw std::runtime_error("error");
  }
  return res;
}

void save_set_string_in_file(std::set<std::string> ls, std::string output_path) {
  std::ofstream fs { output_path, std::fstream::out };
  if (fs.is_open()) {
    for (const auto &s: ls) {
      fs << s << '\n';
    }
  }
  fs.close();
}

std::tuple<std::string, std::string> print_info_song(std::string path_song) {
  const int offset = std::string{"Artist                          : "}.length();
  const std::string cmd = "exiftool " + path_song + " > tmp.txt";
  std::string title {};
  std::string artist {};
  std::string filename {};
  std::system(cmd.data());
  std::ifstream exiftool_output { "tmp.txt", std::fstream::in };
  if (exiftool_output.is_open()) {
    while (exiftool_output.good()) {
      std::string line;
      getline(exiftool_output, line);
      if (line.starts_with("Title")) {
        title = line.substr(offset);
      } else if (line.starts_with("Artist")) {
        artist = line.substr(offset);
      } else if (line.starts_with("File Name")) {
        filename = line.substr(offset);
      }
    }
  } else {
    throw std::runtime_error("error");
  }
  return std::tuple<std::string, std::string>{filename, artist + " - " + title};
}

std::set<std::tuple<std::string, std::string>> get_songs_from_dir(std::string path) {
  std::set<std::tuple<std::string, std::string>> res {};

  for (const auto &file : std::filesystem::directory_iterator(path)) {
    res.insert(print_info_song(file.path()));
  }

  return res;
}

std::set<std::string> get_path_favorite_songs(std::set<std::string> favorite_songs, std::set<std::tuple<std::string, std::string>> repository) {
  std::set<std::string> playlist {};

  for (const auto &fav: favorite_songs) {
    for (const auto &s: repository) {
      if (fav == std::get<1>(s)) {
        playlist.insert(std::get<0>(s));
      }
    }
  }
  return playlist;
}

int main(int argc, char *argv[]) {
  if (argc < 4) {
    std::cout << "main <favorite-song-path> <song-path> <playlist-name>\n";
    return 0;
  }
  try {
    std::set<std::string> favorite_songs = get_songs_from_file(argv[1]);
    std::set<std::tuple<std::string, std::string>> songs = get_songs_from_dir(argv[2]);
    std::set<std::string> playlist = get_path_favorite_songs(favorite_songs, songs);
    save_set_string_in_file(playlist, argv[3]);
  } catch (const std::exception& e) {
    std::cout << e.what() << '\n';
  }

  return 0;
}