package repository

type UserImagesRepository interface {
    SaveUserImage(userId string, imageName string) error
}

type FilesRepository interface {

}
